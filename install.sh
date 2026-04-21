#!/bin/bash

set -e

REPO="TheOjasSingh/kubeaid"
BINARY_NAME="kubeaid"
INSTALL_DIR="/usr/local/bin"

echo "🚀 Installing kubeaid..."

# Detect OS & ARCH

OS="$(uname -s)"
ARCH="$(uname -m)"

if [[ "$OS" != "Linux" ]]; then
echo "❌ Only Linux supported currently"
exit 1
fi

case "$ARCH" in
x86_64) FILE="kubeaid-linux-amd64" ;;
aarch64) FILE="kubeaid-linux-arm64" ;;
*)
echo "❌ Unsupported architecture: $ARCH"
exit 1
;;
esac

URL="https://github.com/$REPO/releases/download/v0.1.0/$FILE"

echo "⬇️ Downloading $FILE..."
curl -fL "$URL" -o "$BINARY_NAME" || {
echo "❌ Download failed. Check release file name."
exit 1
}

echo "🔧 Making executable..."
chmod +x "$BINARY_NAME"

echo "📦 Installing to $INSTALL_DIR..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"

echo "🔍 Verifying installation..."
if command -v kubeaid >/dev/null 2>&1; then
echo "✅ kubeaid installed successfully!"
kubeaid --help || true
else
echo "❌ Installation failed"
exit 1
fi
