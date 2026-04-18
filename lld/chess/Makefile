# Cross-compile the terminal chess app for GitHub Releases.
# Usage: make release && make checksums

BIN := chess
DIST := dist
LDFLAGS := -s -w

.PHONY: release clean build checksums

$(DIST):
	mkdir -p $(DIST)

# Local build for the current OS/architecture.
build:
	CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(BIN) .

# Build all release targets into dist/.
release: $(DIST)
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/chess-linux-amd64 .
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/chess-linux-arm64 .
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/chess-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/chess-darwin-arm64 .
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "$(LDFLAGS)" -o $(DIST)/chess-windows-amd64.exe .

clean:
	rm -rf $(DIST)

# Per-binary .sha256 files plus checksums.txt (run after release).
checksums:
	@set -e; \
	cd $(DIST) || { echo "missing $(DIST); run make release first" >&2; exit 1; }; \
	any=0; \
	for f in chess-linux-amd64 chess-linux-arm64 chess-darwin-amd64 chess-darwin-arm64 chess-windows-amd64.exe; do \
		if [ -f "$$f" ]; then any=1; break; fi; \
	done; \
	if [ "$$any" -eq 0 ]; then echo "no release binaries in $(DIST); run make release first" >&2; exit 1; fi; \
	rm -f *.sha256 checksums.txt; \
	for f in chess-linux-amd64 chess-linux-arm64 chess-darwin-amd64 chess-darwin-arm64 chess-windows-amd64.exe; do \
		test -f "$$f" || continue; \
		if command -v sha256sum >/dev/null 2>&1; then \
			sha256sum "$$f" > "$$f.sha256"; \
		else \
			shasum -a 256 "$$f" | awk -v fn="$$f" '{print $$1 "  " fn}' > "$$f.sha256"; \
		fi; \
	done; \
	cat chess-*.sha256 > checksums.txt; \
	echo "wrote $(DIST)/*.sha256 and $(DIST)/checksums.txt"
