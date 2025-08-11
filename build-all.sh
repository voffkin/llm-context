#!/bin/bash
set -e

OUTPUT_DIR="dist"
mkdir -p "$OUTPUT_DIR"

# Список целей (GOOS/GOARCH)
TARGETS=(
  "linux/amd64"
  "linux/arm64"
  "windows/amd64"
  "windows/arm64"
  "darwin/amd64"
  "darwin/arm64"
)

for target in "${TARGETS[@]}"; do
  GOOS=${target%/*}
  GOARCH=${target#*/}
  
  BIN_NAME="llm-context-${GOOS}-${GOARCH}"
  if [ "$GOOS" = "windows" ]; then
    BIN_NAME="$BIN_NAME.exe"
  fi
  
  echo "==> Building $BIN_NAME"
  GOOS=$GOOS GOARCH=$GOARCH go build -o "$OUTPUT_DIR/$BIN_NAME" ./cmd/llm-context
done

echo "✅ All binaries are in $OUTPUT_DIR/"
