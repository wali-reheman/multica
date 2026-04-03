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
    frontend_child: std::sync::Mutex<Option<std::process::Child>>,
}

impl SidecarState {
    pub fn new(port: u16) -> Self {
        Self {
            port,
            child: std::sync::Mutex::new(None),
            running: AtomicBool::new(false),
            restart_count: AtomicU32::new(0),
            health_handle: std::sync::Mutex::new(None),
            frontend_child: std::sync::Mutex::new(None),
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

/// Wait for a server to become healthy by polling a URL.
pub async fn wait_for_healthy(port: u16) -> Result<(), String> {
    wait_for_url(&format!("http://localhost:{}/health", port)).await
}

/// Wait for the frontend server to respond.
pub async fn wait_for_frontend(port: u16) -> Result<(), String> {
    wait_for_url(&format!("http://localhost:{}/", port)).await
}

async fn wait_for_url(url: &str) -> Result<(), String> {
    let client = reqwest::Client::builder()
        .timeout(std::time::Duration::from_secs(HEALTH_CHECK_TIMEOUT_SECS))
        .build()
        .map_err(|e| format!("failed to create HTTP client: {}", e))?;

    let deadline =
        tokio::time::Instant::now() + std::time::Duration::from_secs(STARTUP_TIMEOUT_SECS);

    while tokio::time::Instant::now() < deadline {
        match client.get(url).send().await {
            Ok(resp) if resp.status().is_success() || resp.status().is_redirection() => {
                log::info!("server healthy: {}", url);
                return Ok(());
            }
            _ => {
                tokio::time::sleep(std::time::Duration::from_millis(500)).await;
            }
        }
    }

    Err(format!(
        "server at {} failed to become healthy within {}s",
        url, STARTUP_TIMEOUT_SECS
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

/// Start the Next.js standalone frontend server.
pub fn start_frontend(app: &AppHandle, frontend_port: u16, api_port: u16) -> Result<(), String> {
    let state = app.state::<Arc<SidecarState>>();

    // Resolve the frontend server.js path.
    // In development: use the web app's .next/standalone/ directory.
    // In production: use the bundled resources/frontend/ directory.
    let resource_dir = app
        .path()
        .resource_dir()
        .map_err(|e| format!("failed to resolve resource dir: {}", e))?;

    let standalone_dir = resource_dir.join("frontend");
    let server_js = standalone_dir.join("server.js");

    // Fall back to the development path
    let (server_js, working_dir) = if server_js.exists() {
        (server_js, standalone_dir)
    } else {
        // Development: look for the standalone build relative to the project
        let dev_standalone = std::env::current_dir()
            .unwrap_or_default()
            .join("../../web/.next/standalone/apps/web");
        let dev_server_js = dev_standalone.join("server.js");
        if dev_server_js.exists() {
            (dev_server_js, dev_standalone)
        } else {
            return Err(
                "Next.js standalone server not found. Run `NEXT_OUTPUT=standalone pnpm build` in apps/web/ first.".into()
            );
        }
    };

    log::info!("starting frontend server on port {} from {:?}", frontend_port, server_js);

    let child = std::process::Command::new("node")
        .arg(&server_js)
        .env("PORT", frontend_port.to_string())
        .env("HOSTNAME", "localhost")
        .env("NEXT_PUBLIC_API_URL", format!("http://localhost:{}", api_port))
        .env("NEXT_PUBLIC_WS_URL", format!("ws://localhost:{}/ws", api_port))
        .current_dir(&working_dir)
        .stdout(std::process::Stdio::piped())
        .stderr(std::process::Stdio::piped())
        .spawn()
        .map_err(|e| format!("failed to start frontend server (is Node.js installed?): {}", e))?;

    *state.frontend_child.lock().unwrap() = Some(child);

    Ok(())
}

/// Gracefully shut down all managed processes.
pub fn shutdown_sidecar(app: &AppHandle) {
    let state = app.state::<Arc<SidecarState>>();

    // Stop health monitor
    if let Some(handle) = state.health_handle.lock().unwrap().take() {
        handle.abort();
    }

    // Kill the Go sidecar child process
    if let Some(child) = state.child.lock().unwrap().take() {
        log::info!("shutting down Go server");
        let _ = child.kill();
    }

    // Kill the frontend server process
    if let Some(mut child) = state.frontend_child.lock().unwrap().take() {
        log::info!("shutting down frontend server");
        let _ = child.kill();
    }

    state.running.store(false, Ordering::SeqCst);
}
