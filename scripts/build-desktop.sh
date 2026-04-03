#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "$0")/.." && pwd)"
BINARIES_DIR="$REPO_ROOT/apps/desktop/src-tauri/binaries"

# Determine the Rust target triple for the current platform
get_target_triple() {
    local arch os
    arch="$(uname -m)"
    os="$(uname -s)"

    case "$arch" in
        arm64|aarch64) arch="aarch64" ;;
        x86_64)        arch="x86_64" ;;
        *)             echo "Unsupported architecture: $arch" >&2; exit 1 ;;
    esac

    case "$os" in
        Darwin) echo "${arch}-apple-darwin" ;;
        Linux)  echo "${arch}-unknown-linux-gnu" ;;
        *)      echo "Unsupported OS: $os" >&2; exit 1 ;;
    esac
}

TARGET_TRIPLE="$(get_target_triple)"

echo "==> Building Go server for sidecar (target: $TARGET_TRIPLE)"
mkdir -p "$BINARIES_DIR"

# Build Go server binary
cd "$REPO_ROOT/server"
CGO_ENABLED=0 go build -o "$BINARIES_DIR/multica-server-${TARGET_TRIPLE}" ./cmd/server

echo "==> Sidecar binary ready: multica-server-${TARGET_TRIPLE}"

echo "==> Building Next.js frontend (standalone mode)"
cd "$REPO_ROOT/apps/web"
NEXT_OUTPUT=standalone npx next build

echo "==> Next.js standalone build complete"

echo "==> Building desktop app..."
cd "$REPO_ROOT/apps/desktop"
cargo tauri build

echo "==> Desktop build complete!"
