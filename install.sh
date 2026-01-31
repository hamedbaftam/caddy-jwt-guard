#!/bin/bash
set -e

echo "🔧 Installing Go..."
if ! command -v go &>/dev/null; then
  curl -LO https://go.dev/dl/go1.22.1.linux-amd64.tar.gz
  sudo tar -C /usr/local -xzf go1.22.1.linux-amd64.tar.gz
fi

export PATH=$PATH:/usr/local/go/bin

echo "📦 Installing xcaddy..."
go install github.com/caddyserver/xcaddy/cmd/xcaddy@latest
export PATH=$PATH:$(go env GOPATH)/bin

echo "🔨 Building Caddy with JWT Guard..."
xcaddy build --with github.com/YOUR_GITHUB_USERNAME/caddy-jwt-guard=./

sudo mv caddy /usr/local/bin/caddy-jwt
sudo chmod +x /usr/local/bin/caddy-jwt

echo "⚙️ Installing configs..."
sudo mkdir -p /etc/caddy
sudo cp Caddyfile.example /etc/caddy/Caddyfile

echo "🧩 Installing systemd service..."
sudo cp systemd/caddy-jwt.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable caddy-jwt
sudo systemctl restart caddy-jwt

echo "✅ Installation complete!"