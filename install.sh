#!/bin/bash

set -euo pipefail

REPO="estevesegura/JanusKnock"
VERSION="0.0.1"

# Detect OS and Architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64|aarch64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Determine binary name
BINARY="janusknock_${VERSION}_${OS}_${ARCH}"
if [[ "$OS" == "windows" ]]; then
    BINARY+=".exe"
fi

URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY"

# Download the binary
echo "Downloading $BINARY from $URL..."
if ! curl -L "$URL" -o /tmp/$BINARY; then
    echo "Error: Failed to download $BINARY"
    exit 1
fi

# Make it executable
chmod +x /tmp/$BINARY

# Move binary to /usr/local/bin
echo "Installing $BINARY to /usr/local/bin..."
if ! sudo mv /tmp/$BINARY /usr/local/bin/janusknock; then
    echo "Error: Failed to move $BINARY to /usr/local/bin"
    exit 1
fi

echo "Installation complete! You can now run 'janusknock'."
