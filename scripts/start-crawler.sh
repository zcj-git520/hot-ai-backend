#!/bin/bash
set -e

cd "$(dirname "$0")/.."
ROOT_DIR=$(pwd)

mkdir -p logs bin deploy/.pids

if [ ! -x bin/crawler-svc ] || [ apps/crawler-svc/main.go -nt bin/crawler-svc ]; then
  echo "[build] go build crawler-svc ..."
  go build -o bin/crawler-svc ./apps/crawler-svc
fi

if pgrep -f "bin/crawler-svc" > /dev/null; then
  echo "[skip] crawler-svc 已在运行 (PID: $(pgrep -f bin/crawler-svc | tr '\n' ' '))"
  exit 0
fi

nohup bin/crawler-svc \
  -f apps/crawler-svc/etc/crawler-svc.yaml \
  > logs/crawler.log 2>&1 &

PID=$!
echo $PID > deploy/.pids/crawler-svc.pid

sleep 1
if kill -0 "$PID" 2>/dev/null; then
  echo "[ok] crawler-svc 已后台启动 PID=$PID，日志: $ROOT_DIR/logs/crawler.log"
else
  echo "[fail] crawler-svc 启动失败，查看日志: $ROOT_DIR/logs/crawler.log"
  tail -n 30 logs/crawler.log
  exit 1
fi