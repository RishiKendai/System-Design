# chess CLI

Source code for this app lives under the System-Design monorepo: [`lld/chess`](https://github.com/RishiKendai/System-Design/tree/main/lld/chess) (same directory as `main.go`).

This is a **two-player chess game** you play **in the terminal**. The UI is built with [**Bubble Tea**](https://github.com/charmbracelet/bubbletea) so it feels **interactive and polished**—full-screen layout, text inputs, and scrollable move history instead of a bare stream of prompts. Two people share the keyboard: you enter names, player 1 picks White or Black (player 2 takes the other side), then take turns entering moves until checkmate, stalemate, or draw. A small engine under the hood handles the board, piece rules, and special cases (castling, en passant, pawn promotion, and the rest).

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

### macOS: “Apple could not verify…” (Gatekeeper)

Prebuilt binaries from GitHub are **not Apple-notarized**, so the first time you open one, macOS may say it cannot verify the app is free of malware. That is normal for unsigned open-source builds; only proceed if you trust this repo and the release you downloaded.

**If you cannot run the binary** (blocked, “damaged,” or Terminal refuses to execute it), remove the download quarantine flag, then `chmod +x` and run again. Adjust the path if the file is not in Downloads; on Intel Macs use `chess-darwin-amd64` instead of `chess-darwin-arm64`.

```bash
xattr -dr com.apple.quarantine ~/Downloads/chess-darwin-arm64
```

**Ways to run it anyway (pick one):**

1. **Finder:** Control-click (or right-click) `chess-darwin-arm64` → **Open** → confirm **Open** on the dialog (not a double-click the first time).
2. **System Settings:** **Privacy & Security** → after you try to open the file once, look for a message about it being blocked and click **Open Anyway**.
3. **Terminal** (good for this app anyway, since it is a full-screen TUI):

   ```bash
   chmod +x /path/to/chess-darwin-arm64
   /path/to/chess-darwin-arm64
   ```

4. If the browser added a **quarantine** flag and Terminal still refuses, remove it (only for files you trust):

   ```bash
   xattr -dr com.apple.quarantine /path/to/chess-darwin-arm64
   ```

Long term, a maintainer can distribute a **signed and notarized** macOS build to avoid this prompt; that requires an Apple Developer Program subscription and a release pipeline that submits binaries to Apple.

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
