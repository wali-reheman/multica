mod commands;
mod sidecar;
mod tray;

use std::sync::Arc;

use sidecar::SidecarState;
use tauri::Emitter;

pub fn run() {
    let port = sidecar::find_free_port();

    tauri::Builder::default()
        .plugin(tauri_plugin_shell::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(tauri_plugin_notification::init())
        .plugin(tauri_plugin_process::init())
        .plugin(tauri_plugin_opener::init())
        .plugin(
            tauri_plugin_log::Builder::new()
                .target(tauri_plugin_log::Target::new(
                    tauri_plugin_log::TargetKind::Stdout,
                ))
                .target(tauri_plugin_log::Target::new(
                    tauri_plugin_log::TargetKind::LogDir {
                        file_name: Some("multica".into()),
                    },
                ))
                .build(),
        )
        .plugin(tauri_plugin_updater::Builder::new().build())
        .manage(Arc::new(SidecarState::new(port)))
        .invoke_handler(tauri::generate_handler![
            commands::get_server_info,
            commands::pick_folder,
            commands::send_notification,
            commands::restart_server,
        ])
        .setup(move |app| {
            let app_handle = app.handle().clone();

            // Start the Go server sidecar
            sidecar::start_sidecar(&app_handle)
                .map_err(|e| Box::new(std::io::Error::new(std::io::ErrorKind::Other, e)))?;

            // Wait for the server to be healthy before showing the window
            let setup_handle = app_handle.clone();
            tauri::async_runtime::spawn(async move {
                match sidecar::wait_for_healthy(port).await {
                    Ok(()) => {
                        log::info!("server ready, creating main window");
                        create_main_window(&setup_handle, port);
                        sidecar::start_health_monitor(&setup_handle);
                    }
                    Err(e) => {
                        log::error!("server failed to start: {}", e);
                        // Show window anyway so user can see the error
                        create_main_window(&setup_handle, port);
                        let _ = setup_handle.emit("sidecar-crashed", ());
                    }
                }
            });

            // Create system tray
            if let Err(e) = tray::create_tray(&app_handle) {
                log::warn!("failed to create system tray: {}", e);
            }

            Ok(())
        })
        .on_window_event(|window, event| {
            if let tauri::WindowEvent::CloseRequested { api, .. } = event {
                // Hide instead of close — keep running in tray
                let _ = window.hide();
                api.prevent_close();
            }
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

fn create_main_window(app: &tauri::AppHandle, port: u16) {
    let init_script = format!(
        r#"
        window.__MULTICA_API_URL__ = "http://localhost:{}";
        window.__MULTICA_WS_URL__ = "ws://localhost:{}/ws";
        window.__TAURI_DESKTOP__ = true;
        "#,
        port, port
    );

    let builder = tauri::WebviewWindowBuilder::new(app, "main", tauri::WebviewUrl::App("index.html".into()))
        .title("Multica")
        .inner_size(1200.0, 800.0)
        .min_inner_size(900.0, 600.0)
        .initialization_script(&init_script);

    #[cfg(target_os = "macos")]
    let builder = builder
        .title_bar_style(tauri::TitleBarStyle::Overlay)
        .hidden_title(true);

    match builder.build() {
        Ok(_) => log::info!("main window created"),
        Err(e) => log::error!("failed to create main window: {}", e),
    }
}
