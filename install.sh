#!/bin/bash
set -e

REPO="your-username/JanusKnock"
VERSION="v1.0.0"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

if [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
    ARCH="arm64"
else
    echo "Unsupported architecture: $ARCH"
    exit 1
fi

BINARY="janusknock_${VERSION}_${OS}_${ARCH}"
if [[ "$OS" == "windows" ]]; then
    BINARY+=".exe"
fi

URL="https://github.com/$REPO/releases/download/$VERSION/$BINARY"

echo "Downloading $BINARY from $URL..."
curl -L "$URL" -o /tmp/$BINARY
chmod +x /tmp/$BINARY

echo "Installing $BINARY to /usr/local/bin..."
sudo mv /tmp/$BINARY /usr/local/bin/janusknock

echo "Installation complete! You can now run 'janusknock'."
