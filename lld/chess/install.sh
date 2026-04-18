#!/usr/bin/env bash
# Install the terminal chess binary from GitHub Releases (Linux and macOS).
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/RishiKendai/System-Design/main/lld/chess/install.sh | bash
# Or with explicit repo (e.g. a fork):
#   CHESS_INSTALL_GITHUB_OWNER=you CHESS_INSTALL_GITHUB_REPO=your-repo bash install.sh
#
# Environment:
#   CHESS_INSTALL_GITHUB_OWNER   default: RishiKendai (GitHub user/org for Releases)
#   CHESS_INSTALL_GITHUB_REPO    default: System-Design (release assets live on this repo)
#   CHESS_INSTALL_VERSION        empty = latest; else tag, e.g. v1.0.0 (releases/download/<tag>/...)
#   CHESS_INSTALL_PREFIX         install directory (default: $HOME/.local/bin)
#   CHESS_INSTALL_BINARY_NAME    installed filename (default: chess)

set -euo pipefail

OWNER="${CHESS_INSTALL_GITHUB_OWNER:-RishiKendai}"
REPO="${CHESS_INSTALL_GITHUB_REPO:-System-Design}"
VERSION="${CHESS_INSTALL_VERSION:-}"
PREFIX="${CHESS_INSTALL_PREFIX:-$HOME/.local/bin}"
BINARY_NAME="${CHESS_INSTALL_BINARY_NAME:-chess}"

case "$(uname -s)" in
Linux*)
	goos=linux
	;;
Darwin*)
	goos=darwin
	;;
*)
	echo "install.sh: This script supports Linux and macOS only." >&2
	echo "On Windows, download chess-windows-amd64.exe from GitHub Releases." >&2
	exit 1
	;;
esac

case "$(uname -m)" in
x86_64 | amd64)
	goarch=amd64
	;;
aarch64 | arm64)
	goarch=arm64
	;;
*)
	echo "install.sh: Unsupported CPU architecture: $(uname -m)" >&2
	exit 1
	;;
esac

asset="chess-${goos}-${goarch}"
if [ -n "${VERSION}" ]; then
	base="https://github.com/${OWNER}/${REPO}/releases/download/${VERSION}"
else
	base="https://github.com/${OWNER}/${REPO}/releases/latest/download"
fi

tmp="$(mktemp -d)"
cleanup() {
	rm -rf "${tmp}"
}
trap cleanup EXIT

echo "Downloading ${base}/${asset} ..."
curl -fsSL "${base}/${asset}" -o "${tmp}/${asset}"

if curl -fsSL "${base}/${asset}.sha256" -o "${tmp}/${asset}.sha256" 2>/dev/null; then
	echo "Verifying checksum ..."
	(
		cd "${tmp}"
		if command -v sha256sum >/dev/null 2>&1; then
			sha256sum -c "${asset}.sha256"
		else
			shasum -a 256 -c "${asset}.sha256"
		fi
	)
else
	echo "install.sh: No .sha256 sidecar found; skipping checksum verification." >&2
fi

chmod +x "${tmp}/${asset}"
mkdir -p "${PREFIX}"

if ! cp "${tmp}/${asset}" "${PREFIX}/${BINARY_NAME}" 2>/dev/null; then
	echo "install.sh: Could not write to ${PREFIX}/${BINARY_NAME}" >&2
	echo "Try: CHESS_INSTALL_PREFIX=/usr/local/bin sudo -E env \"PATH=\$PATH\" bash install.sh" >&2
	echo "Or: mkdir -p \"${PREFIX}\" && fix permissions, then re-run." >&2
	exit 1
fi

echo "Installed: ${PREFIX}/${BINARY_NAME}"
if ! command -v "${BINARY_NAME}" >/dev/null 2>&1; then
	case ":${PATH}:" in
	*":${PREFIX}:"*) ;;
	*)
		echo "Add ${PREFIX} to PATH, for example:" >&2
		echo "  export PATH=\"${PREFIX}:\$PATH\"" >&2
		;;
	esac
fi
