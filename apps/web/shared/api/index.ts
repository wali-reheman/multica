import { createLogger } from "@/shared/logger";
import { ApiClient } from "./client";

export { ApiClient } from "./client";
export type { LoginResponse } from "./client";
export { WSClient } from "./ws-client";

// In Tauri desktop mode, the sidecar port is injected at runtime via initialization script.
// Falls back to build-time env var, then empty string (same-origin).
const API_BASE_URL =
  (typeof window !== "undefined" && (window as any).__MULTICA_API_URL__) ||
  process.env.NEXT_PUBLIC_API_URL ||
  "";

export const api = new ApiClient(API_BASE_URL, { logger: createLogger("api") });

// Initialize token from localStorage on load
if (typeof window !== "undefined") {
  const token = localStorage.getItem("multica_token");
  if (token) {
    api.setToken(token);
  }
  const wsId = localStorage.getItem("multica_workspace_id");
  if (wsId) {
    api.setWorkspaceId(wsId);
  }

}
