---
plan: "05-02"
phase: "05-polish-production-deployment"
status: complete
completed: 2026-04-03
self_check: PASSED
---

# Plan 05-02: Build Pipeline + Deployment Runbook — SUMMARY

## What was built

Wired the production build pipeline and wrote the Raspberry Pi deployment runbook, satisfying DEPL-01–04.

## Key files created/modified

### Modified
- `Makefile` — `build` target updated (npm run build before go build), `build-pi` target added (arm64 cross-compile), `.PHONY` updated
- `.planning/REQUIREMENTS.md` — UI-04 updated to reflect Lottie decision (D-01 override); stale "Lottie for sidebar animation" Out-of-Scope row removed

### Created
- `docs/deployment.md` — Complete Raspberry Pi deployment runbook (7 steps + troubleshooting + env var reference)

## Verification results

- **DEPL-01**: `make build` succeeded — 21MB `flos-library` binary produced with embedded React app
- **DEPL-02**: systemd `floslib.service` unit documented with `EnvironmentFile`, `WorkingDirectory=/opt/floslib`, `Restart=on-failure`
- **DEPL-03**: Caddy `Caddyfile` with `reverse_proxy localhost:8081` + automatic Let's Encrypt TLS documented
- **DEPL-04**: `GOOGLE_BOOKS_API_KEY` absent from `frontend/dist/` JS bundle — audit passed

## Runbook coverage

- Step 1: `make build-pi` cross-compile (CGO_ENABLED=0 GOOS=linux GOARCH=arm64)
- Step 2: Pi setup — PostgreSQL, app directory, .env file
- Step 3: Database migrations via SSH tunnel
- Step 4: scp binary + data/covers/ transfer
- Step 5: systemd service install + enable
- Step 6: Caddy HTTPS configuration
- Step 7: Subsequent deployment update flow (~1s downtime)
- Troubleshooting table + environment variable reference

## Notes

- `make build-pi` on Windows: `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build` works via GNU Make (Git Bash/MSYS2 sh). For manual PowerShell cross-compile, env vars must be set separately
- Vite eval warning from lottie-web is expected — it's inside the library's expression evaluator, not a security issue in this app
