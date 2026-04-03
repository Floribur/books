# Flo's Library — Deployment Runbook

> **Target:** Raspberry Pi 4/5 running Raspberry Pi OS (64-bit / arm64)
> **Approach:** Single Go binary serving embedded React frontend + PostgreSQL as DB
> **Phase:** This runbook covers production deployment steps. Run these manually when ready to go live.

---

## Architecture Overview

```
Internet → Caddy (HTTPS, port 443) → flos-library binary (HTTP, port 8081)
                                          ↓
                                     PostgreSQL (local, port 5432)
                                          ↓
                                     data/covers/ (local directory, NOT embedded)
```

**What the binary contains (embedded at compile time):**
- Entire React frontend (`frontend/dist/`) — served as a single-page app
- All static assets (JS, CSS, fonts)

**What stays on disk (NOT embedded):**
- Book cover images: `data/covers/` — must exist alongside the binary
- Environment file: `.env` — secrets; never commit to git

---

## Step 1: Cross-Compile Binary for Raspberry Pi

On your **development machine** (Windows/Mac/Linux), from the project root:

```bash
make build-pi
```

This runs:
1. `cd frontend && npm run build` — bundles React into `frontend/dist/`
2. `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o flos-library-linux-arm64 ./cmd/server`

Output: `flos-library-linux-arm64` in the project root.

> **Why `CGO_ENABLED=0`:** Go's cross-compilation requires disabling cgo when no cross-toolchain is
> installed. The app uses `pgx/v5` (pure Go) and has no cgo dependencies — this is safe.

---

## Step 2: Prepare the Pi

SSH into your Raspberry Pi and set up the environment:

```bash
# Install PostgreSQL (if not already installed)
sudo apt update && sudo apt install -y postgresql

# Create database and user
sudo -u postgres psql -c "CREATE USER floslib WITH PASSWORD 'your_strong_password';"
sudo -u postgres psql -c "CREATE DATABASE floslib OWNER floslib;"

# Create application directory
sudo mkdir -p /opt/floslib
sudo mkdir -p /opt/floslib/data/covers

# Create .env file with secrets (edit values before saving)
sudo tee /opt/floslib/.env > /dev/null <<'EOF'
DATABASE_URL=postgres://floslib:your_strong_password@localhost:5432/floslib?sslmode=disable
GOOGLE_BOOKS_API_KEY=your_google_books_api_key_here
PORT=8081
EOF

sudo chmod 600 /opt/floslib/.env
```

> **Security:** `GOOGLE_BOOKS_API_KEY` lives only in `/opt/floslib/.env` on the Pi (backend process).
> It is NEVER embedded in the React/JS bundle — Vite only exposes `VITE_`-prefixed variables to the browser.

---

## Step 3: Run Database Migrations on the Pi

From your **development machine** (requires `golang-migrate` installed):

```bash
# Forward Pi's PostgreSQL port locally, then run migrations
ssh -L 5433:localhost:5432 pi@<pi-ip-address> -N &
migrate -path ./migrations \
  -database "postgres://floslib:your_strong_password@localhost:5433/floslib?sslmode=disable" up
```

Or SSH into the Pi, copy `migrations/` folder, and run migrate there directly.

---

## Step 4: Transfer Binary and Assets to the Pi

From your **development machine**:

```bash
# Transfer binary
scp flos-library-linux-arm64 pi@<pi-ip-address>:/opt/floslib/flos-library

# Transfer cover images (if first deployment; subsequent syncs happen via the app)
scp -r data/covers/ pi@<pi-ip-address>:/opt/floslib/data/

# Make binary executable
ssh pi@<pi-ip-address> "chmod +x /opt/floslib/flos-library"
```

---

## Step 5: Install as systemd Service

SSH into the Pi and create the service unit file:

```bash
sudo tee /etc/systemd/system/floslib.service > /dev/null <<'EOF'
[Unit]
Description=Flo's Library — book showcase server
After=network.target postgresql.service
Requires=postgresql.service

[Service]
Type=simple
User=pi
WorkingDirectory=/opt/floslib
EnvironmentFile=/opt/floslib/.env
ExecStart=/opt/floslib/flos-library
Restart=on-failure
RestartSec=5s

# Security hardening
NoNewPrivileges=true
PrivateTmp=true

[Install]
WantedBy=multi-user.target
EOF

# Enable and start
sudo systemctl daemon-reload
sudo systemctl enable floslib
sudo systemctl start floslib

# Verify running
sudo systemctl status floslib
```

> **Note:** `WorkingDirectory=/opt/floslib` is important — the binary reads `data/covers/` relative
> to its working directory. If this is wrong, cover images will return 404.

---

## Step 6: Configure Caddy for HTTPS

Caddy handles TLS automatically via Let's Encrypt (zero-config HTTPS).

Install Caddy on the Pi:
```bash
sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
sudo apt update && sudo apt install caddy
```

Create `/etc/caddy/Caddyfile`:
```
your-domain.com {
    reverse_proxy localhost:8081
}
```

Start and enable Caddy:
```bash
sudo systemctl enable caddy
sudo systemctl start caddy
```

Caddy automatically obtains and renews TLS certificates from Let's Encrypt. No certbot required.

> **DNS requirement:** Your domain must point to the Pi's public IP before Caddy can obtain a certificate.

---

## Step 7: Subsequent Deployments (Update Flow)

When you have a new version to deploy:

```bash
# 1. Build new binary on dev machine
make build-pi

# 2. Transfer binary
scp flos-library-linux-arm64 pi@<pi-ip-address>:/opt/floslib/flos-library

# 3. Restart service
ssh pi@<pi-ip-address> "sudo systemctl restart floslib"

# 4. Verify
ssh pi@<pi-ip-address> "sudo systemctl status floslib"
```

Total downtime: ~1 second (systemd restart).

---

## Troubleshooting

| Symptom | Cause | Fix |
|---------|-------|-----|
| `dial tcp 127.0.0.1:5432: connect: connection refused` | PostgreSQL not running | `sudo systemctl start postgresql` |
| `data/covers/*.jpg` returns 404 | Wrong WorkingDirectory or missing covers dir | Check `WorkingDirectory` in systemd unit; verify `/opt/floslib/data/covers/` exists |
| Caddy TLS: `no such host` | DNS not propagated yet | Wait for DNS propagation; check with `dig your-domain.com` |
| Binary exits immediately | Check logs: `sudo journalctl -u floslib -n 50` | Usually missing DATABASE_URL in .env |
| Old React version served | Go embed cached stale build | Always run `make build-pi` (not `go build` alone) — npm build must precede go build |

---

## Environment Variables Reference

| Variable | Required | Description | Example |
|----------|----------|-------------|---------|
| `DATABASE_URL` | YES | PostgreSQL connection string | `postgres://floslib:pass@localhost:5432/floslib?sslmode=disable` |
| `GOOGLE_BOOKS_API_KEY` | YES (for enrichment) | Google Books API key — backend only | `AIzaSy...` |
| `PORT` | NO (default: 8081) | HTTP port the Go server listens on | `8081` |

> **Security reminder:** `GOOGLE_BOOKS_API_KEY` must NEVER be prefixed with `VITE_`. Vite only exposes
> `VITE_*` variables to the browser bundle. Verify after any bundle change:
> `grep -r "GOOGLE_BOOKS_API_KEY" frontend/dist/` — must return nothing.

---

*Last updated: Phase 5 — Polish & Production Deployment*
*Target: Raspberry Pi 4/5 (arm64 / ARMv8)*
*Binary built with: `CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build`*
