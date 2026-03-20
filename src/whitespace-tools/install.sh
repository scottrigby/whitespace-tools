#!/usr/bin/env bash
set -euo pipefail

VERSION="${VERSION:-"latest"}"
REPO="scottrigby/whitespace-tools"

# Resolve 'latest' to actual version tag
if [ "${VERSION}" = "latest" ]; then
    VERSION=$(curl -fsSL "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name"' | sed -E 's/.*"tag_name": *"([^"]+)".*/\1/')
fi

# Strip leading 'v' for asset filename
VERSION_NUM="${VERSION#v}"

# Map OS
OS="$(uname -s)"
case "${OS}" in
    Linux)  OS="linux" ;;
    Darwin) OS="darwin" ;;
    *)
        echo "Unsupported OS: ${OS}"
        exit 1
        ;;
esac

# Map architecture
ARCH="$(uname -m)"
case "${ARCH}" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    *)
        echo "Unsupported architecture: ${ARCH}"
        exit 1
        ;;
esac
ASSET="whitespace-tools_${VERSION_NUM}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/${REPO}/releases/download/${VERSION}/${ASSET}"

echo "Installing whitespace-tools ${VERSION} (${OS}/${ARCH})..."

TMP=$(mktemp -d)
trap 'rm -rf "${TMP}"' EXIT

curl -fsSL "${URL}" -o "${TMP}/${ASSET}"
tar -xzf "${TMP}/${ASSET}" -C "${TMP}"

install -m 0755 "${TMP}/newline" /usr/local/bin/newline
install -m 0755 "${TMP}/trailingspace" /usr/local/bin/trailingspace

echo "Installed:"
echo "  $(newline --version)"
echo "  $(trailingspace --version)"
