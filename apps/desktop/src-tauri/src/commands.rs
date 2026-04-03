use std::sync::Arc;

use serde::Serialize;
use tauri::{AppHandle, Manager};

use crate::sidecar::SidecarState;

#[derive(Serialize)]
pub struct ServerInfo {
    pub port: u16,
    pub api_url: String,
    pub ws_url: String,
}

/// Returns the local server connection info.
#[tauri::command]
pub fn get_server_info(app: AppHandle) -> ServerInfo {
    let state = app.state::<Arc<SidecarState>>();
    let port = state.port;
    ServerInfo {
        port,
        api_url: format!("http://localhost:{}", port),
        ws_url: format!("ws://localhost:{}/ws", port),
    }
}

/// Opens a native file picker dialog and returns the selected path.
#[tauri::command]
pub async fn pick_folder(app: AppHandle) -> Result<Option<String>, String> {
    use tauri_plugin_dialog::DialogExt;

    let folder = app.dialog().file().blocking_pick_folder();
    Ok(folder.map(|p| p.to_string()))
}

/// Sends a native OS notification.
#[tauri::command]
pub async fn send_notification(
    app: AppHandle,
    title: String,
    body: String,
) -> Result<(), String> {
    use tauri_plugin_notification::NotificationExt;

    app.notification()
        .builder()
        .title(&title)
        .body(&body)
        .show()
        .map_err(|e| format!("notification error: {}", e))
}

/// Restarts the sidecar server process.
#[tauri::command]
pub fn restart_server(app: AppHandle) -> Result<(), String> {
    crate::sidecar::shutdown_sidecar(&app);
    crate::sidecar::start_sidecar(&app)
}
