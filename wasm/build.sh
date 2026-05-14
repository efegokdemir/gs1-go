#!/bin/bash
set -euo pipefail

# Build the GS1 WASM module.
# Run from the gs1/ directory: bash wasm/build.sh

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
GS1_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$GS1_DIR"

echo "Building gs1.wasm..."
GOOS=js GOARCH=wasm go build -ldflags="-s -w" -o wasm/gs1.wasm ./cmd/gs1-wasm/

echo "Copying wasm_exec.js from Go SDK..."
GOROOT="$(go env GOROOT)"
WASM_EXEC=""
for candidate in "$GOROOT/misc/wasm/wasm_exec.js" "$GOROOT/lib/wasm/wasm_exec.js"; do
    if [ -f "$candidate" ]; then
        WASM_EXEC="$candidate"
        break
    fi
done
if [ -z "$WASM_EXEC" ]; then
    echo "Error: wasm_exec.js not found in Go SDK" >&2
    exit 1
fi
cp "$WASM_EXEC" wasm/wasm_exec.js

# Copy artifacts into example/ so it can be served standalone.
cp wasm/gs1.wasm wasm/example/gs1.wasm
cp wasm/wasm_exec.js wasm/example/wasm_exec.js

SIZE=$(ls -lh wasm/gs1.wasm | awk '{print $5}')
echo "Build complete: wasm/gs1.wasm ($SIZE)"
echo "Example ready: python3 -m http.server 8080 --directory wasm/example/"
