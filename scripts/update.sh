#!/bin/bash
set -euo pipefail
PROJECT_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$PROJECT_DIR"
echo "[$(date)] Deploying..." && git fetch origin main
LOCAL=$(git rev-parse HEAD); REMOTE=$(git rev-parse origin/main)
if [ "$LOCAL" = "$REMOTE" ]; then echo "Up to date."; exit 0; fi
git pull origin main && docker compose build --no-cache && docker compose up -d --force-recreate
sleep 5 && docker compose ps | grep -q "Up" && echo "✅ Deploy OK" || (echo "❌ Deploy failed" && docker compose logs --tail=20 && exit 1)
