# Terminal chess

Source code for this app lives under the System-Design monorepo: [`lld/chess`](https://github.com/RishiKendai/System-Design/tree/main/lld/chess) (same directory as `main.go`).

This is a **two-player chess game** you play **in the terminal**. The UI is built with [**Bubble Tea**](https://github.com/charmbracelet/bubbletea) so it feels **interactive and polished**—full-screen layout, text inputs, and scrollable move history instead of a bare stream of prompts. Two people share the keyboard: you enter names and colors, then take turns entering moves until checkmate, stalemate, or draw. A small engine under the hood handles the board, piece rules, and special cases (castling, en passant, pawn promotion, and the rest).

## What you get

- Two-player chess in the terminal with an alt-screen TUI.
- Move entry by coordinates (e.g. `e2 e4`) or row/column indices; `help` / `?` in-game lists formats.
- Pawn promotion to queen, rook, bishop, or knight.
- Rules and special moves handled in the game layer so the UI stays focused on play.
- Static release binaries for Linux, macOS, and Windows (`CGO_ENABLED=0` cross-builds).

## Requirements

- **Prebuilt binary:** A modern 64-bit terminal (Windows Terminal, macOS Terminal, Linux console or emulator).
- **From source:** [Go](https://go.dev/dl/) **1.21+** (match or exceed the `go` line in [`go.mod`](go.mod)).

## Install from GitHub Releases

Binaries are published on [**Releases** for `RishiKendai/System-Design`](https://github.com/RishiKendai/System-Design/releases). Download the file that matches your OS and CPU.

| OS       | CPU           | Release asset (example)        |
|----------|---------------|--------------------------------|
| Linux    | x86_64        | `chess-linux-amd64`            |
| Linux    | ARM64         | `chess-linux-arm64`            |
| macOS    | Intel         | `chess-darwin-amd64`           |
| macOS    | Apple Silicon | `chess-darwin-arm64`           |
| Windows  | x86_64        | `chess-windows-amd64.exe`      |

### Linux and macOS (manual)

```bash
chmod +x chess-linux-amd64   # or chess-darwin-arm64, etc.
./chess-linux-amd64
```

Install on your PATH (system-wide example):

```bash
sudo mv chess-linux-amd64 /usr/local/bin/chess
chess
```

Or for your user only:

```bash
mkdir -p "$HOME/.local/bin"
mv chess-linux-amd64 "$HOME/.local/bin/chess"
chmod +x "$HOME/.local/bin/chess"
# Put ~/.local/bin on PATH, then:
chess
```

### Optional checksum verification

If you upload `*.sha256` files from `make checksums`:

```bash
sha256sum -c chess-linux-amd64.sha256   # Linux
# macOS:
shasum -a 256 -c chess-darwin-arm64.sha256
```

### One-line install (Linux / macOS)

Review any script before piping it to a shell. [`install.sh`](https://github.com/RishiKendai/System-Design/blob/main/lld/chess/install.sh) pulls the matching binary from [Releases](https://github.com/RishiKendai/System-Design/releases) (upload the `dist/*` outputs from `make release`).

```bash
curl -fsSL https://raw.githubusercontent.com/RishiKendai/System-Design/main/lld/chess/install.sh | bash
```

For another Releases location (e.g. a fork):

```bash
export CHESS_INSTALL_GITHUB_OWNER=YOUR_GITHUB_USER
export CHESS_INSTALL_GITHUB_REPO=YOUR_REPO
curl -fsSL "https://raw.githubusercontent.com/${CHESS_INSTALL_GITHUB_OWNER}/${CHESS_INSTALL_GITHUB_REPO}/main/lld/chess/install.sh" | bash
```

Environment variables (see [`install.sh`](install.sh)):

- `CHESS_INSTALL_VERSION` — unset for **latest**; or a tag like `v1.0.0` for `releases/download/v1.0.0/`.
- `CHESS_INSTALL_PREFIX` — default `$HOME/.local/bin`; `/usr/local/bin` for a shared install if you have permission.
- `CHESS_INSTALL_BINARY_NAME` — default `chess`.

### Windows

1. Download `chess-windows-amd64.exe` from Releases.
2. Run from PowerShell or cmd, for example:

```powershell
.\chess-windows-amd64.exe
```

Rename to `chess.exe` and add the folder to PATH if you want a shorter command.

## Run

With the binary on your `PATH` as `chess`:

```bash
chess
```

Otherwise run it by path (see above). The app uses the terminal’s alternate screen; use `quit` / `q` to exit (see in-app help).

## Develop and build from source

Clone [RishiKendai/System-Design](https://github.com/RishiKendai/System-Design), then:

```bash
git clone https://github.com/RishiKendai/System-Design.git
cd System-Design/lld/chess
go run .
```

Local binary via Make:

```bash
make build
./chess
```

### Cross-compile for releases

```bash
make release    # writes binaries under dist/
make checksums # optional: dist/*.sha256 and dist/checksums.txt
```

Upload the `dist/` artifacts you care about to a [GitHub Release](https://docs.github.com/en/repositories/releasing-projects-on-github/managing-releases-in-a-repository). Names must match [`install.sh`](install.sh): `chess-<os>-<arch>` (and `.exe` on Windows), plus optional `.sha256` sidecars.

### Install with `go install`

Module path ([`go.mod`](go.mod)):

```text
github.com/RishiKendai/System-Design/lld/chess
```

```bash
go install github.com/RishiKendai/System-Design/lld/chess@latest
```

Put `$GOPATH/bin` or `$HOME/go/bin` on your `PATH`. The installed command is typically named **`chess`** (last path segment of the module you install).

## Makefile targets

| Target           | Description                                |
|------------------|--------------------------------------------|
| `make build`     | Build `chess` for the current OS/arch      |
| `make release`   | Cross-compile all release binaries to `dist/` |
| `make checksums` | SHA256 sidecars and `dist/checksums.txt` |
| `make clean`     | Remove `dist/`                             |

## License

See the parent repository or add a `LICENSE` file as appropriate for your project.
