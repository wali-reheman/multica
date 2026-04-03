use std::net::TcpListener;
use std::sync::atomic::{AtomicBool, AtomicU32, Ordering};
use std::sync::Arc;

use tauri::async_runtime::JoinHandle;
use tauri::{AppHandle, Emitter, Manager};
use tauri_plugin_shell::process::CommandChild;
use tauri_plugin_shell::ShellExt;

const MAX_RESTART_ATTEMPTS: u32 = 3;
const HEALTH_CHECK_INTERVAL_SECS: u64 = 5;
const HEALTH_CHECK_TIMEOUT_SECS: u64 = 3;
const STARTUP_TIMEOUT_SECS: u64 = 30;

pub struct SidecarState {
    pub port: u16,
    child: std::sync::Mutex<Option<CommandChild>>,
    running: AtomicBool,
    restart_count: AtomicU32,
    health_handle: std::sync::Mutex<Option<JoinHandle<()>>>,
}

impl SidecarState {
    pub fn new(port: u16) -> Self {
        Self {
            port,
            child: std::sync::Mutex::new(None),
            running: AtomicBool::new(false),
            restart_count: AtomicU32::new(0),
            health_handle: std::sync::Mutex::new(None),
        }
    }
}

/// Find an available TCP port on localhost.
pub fn find_free_port() -> u16 {
    let listener = TcpListener::bind("127.0.0.1:0").expect("failed to bind to find free port");
    listener.local_addr().unwrap().port()
}

/// Start the Go server sidecar process.
pub fn start_sidecar(app: &AppHandle) -> Result<(), String> {
    let state = app.state::<Arc<SidecarState>>();
    let port = state.port;

    log::info!("starting sidecar on port {}", port);

    let sidecar_command = app
        .shell()
        .sidecar("multica-server")
        .map_err(|e| format!("failed to create sidecar command: {}", e))?
        .env("PORT", port.to_string());

    let (mut rx, child) = sidecar_command
        .spawn()
        .map_err(|e| format!("failed to spawn sidecar: {}", e))?;

    *state.child.lock().unwrap() = Some(child);
    state.running.store(true, Ordering::SeqCst);
    state.restart_count.store(0, Ordering::SeqCst);

    // Pipe sidecar stdout/stderr to Tauri's log system
    let app_handle = app.clone();
    tauri::async_runtime::spawn(async move {
        use tauri_plugin_shell::process::CommandEvent;
        while let Some(event) = rx.recv().await {
            match event {
                CommandEvent::Stdout(line) => {
                    let text = String::from_utf8_lossy(&line);
                    log::info!("[server] {}", text.trim());
                }
                CommandEvent::Stderr(line) => {
                    let text = String::from_utf8_lossy(&line);
                    log::warn!("[server] {}", text.trim());
                }
                CommandEvent::Terminated(payload) => {
                    log::error!(
                        "[server] terminated with code {:?}, signal {:?}",
                        payload.code,
                        payload.signal
                    );
                    let state = app_handle.state::<Arc<SidecarState>>();
                    state.running.store(false, Ordering::SeqCst);
                    *state.child.lock().unwrap() = None;

                    // Attempt restart if under the limit
                    let count = state.restart_count.fetch_add(1, Ordering::SeqCst);
                    if count < MAX_RESTART_ATTEMPTS {
                        log::info!(
                            "restarting sidecar (attempt {}/{})",
                            count + 1,
                            MAX_RESTART_ATTEMPTS
                        );
                        tokio::time::sleep(std::time::Duration::from_secs(2)).await;
                        if let Err(e) = start_sidecar(&app_handle) {
                            log::error!("failed to restart sidecar: {}", e);
                        }
                    } else {
                        log::error!("sidecar exceeded max restart attempts, giving up");
                        // Emit event to frontend
                        let _ = app_handle.emit("sidecar-crashed", ());
                    }
                    break;
                }
                _ => {}
            }
        }
    });

    Ok(())
}

/// Wait for the sidecar to become healthy by polling /health.
pub async fn wait_for_healthy(port: u16) -> Result<(), String> {
    let client = reqwest::Client::builder()
        .timeout(std::time::Duration::from_secs(HEALTH_CHECK_TIMEOUT_SECS))
        .build()
        .map_err(|e| format!("failed to create HTTP client: {}", e))?;

    let url = format!("http://localhost:{}/health", port);
    let deadline =
        tokio::time::Instant::now() + std::time::Duration::from_secs(STARTUP_TIMEOUT_SECS);

    while tokio::time::Instant::now() < deadline {
        match client.get(&url).send().await {
            Ok(resp) if resp.status().is_success() => {
                log::info!("sidecar is healthy on port {}", port);
                return Ok(());
            }
            _ => {
                tokio::time::sleep(std::time::Duration::from_millis(500)).await;
            }
        }
    }

    Err(format!(
        "sidecar failed to become healthy within {}s",
        STARTUP_TIMEOUT_SECS
    ))
}

/// Start a background health check loop that monitors the sidecar.
pub fn start_health_monitor(app: &AppHandle) {
    let state = app.state::<Arc<SidecarState>>();
    let port = state.port;
    let app_handle = app.clone();

    let handle = tauri::async_runtime::spawn(async move {
        let client = reqwest::Client::builder()
            .timeout(std::time::Duration::from_secs(HEALTH_CHECK_TIMEOUT_SECS))
            .build()
            .unwrap();

        let url = format!("http://localhost:{}/health", port);
        let mut consecutive_failures: u32 = 0;

        loop {
            tokio::time::sleep(std::time::Duration::from_secs(HEALTH_CHECK_INTERVAL_SECS)).await;

            let state = app_handle.state::<Arc<SidecarState>>();
            if !state.running.load(Ordering::SeqCst) {
                continue;
            }

            match client.get(&url).send().await {
                Ok(resp) if resp.status().is_success() => {
                    consecutive_failures = 0;
                    // Reset restart count on successful health check
                    state.restart_count.store(0, Ordering::SeqCst);
                }
                _ => {
                    consecutive_failures += 1;
                    log::warn!(
                        "sidecar health check failed ({} consecutive)",
                        consecutive_failures
                    );
                    if consecutive_failures >= 3 {
                        log::error!("sidecar unresponsive after 3 health checks, restarting");
                        // Kill and restart
                        shutdown_sidecar(&app_handle);
                        if let Err(e) = start_sidecar(&app_handle) {
                            log::error!("failed to restart unresponsive sidecar: {}", e);
                        }
                        consecutive_failures = 0;
                    }
                }
            }
        }
    });

    *state.health_handle.lock().unwrap() = Some(handle);
}

/// Gracefully shut down the sidecar process.
pub fn shutdown_sidecar(app: &AppHandle) {
    let state = app.state::<Arc<SidecarState>>();

    // Stop health monitor
    if let Some(handle) = state.health_handle.lock().unwrap().take() {
        handle.abort();
    }

    // Kill the sidecar child process
    if let Some(child) = state.child.lock().unwrap().take() {
        log::info!("shutting down sidecar");
        let _ = child.kill();
    }
    state.running.store(false, Ordering::SeqCst);
}
