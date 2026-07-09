#!/bin/bash

cd "$(dirname "$0")/.."

PID_FILE="deploy/.pids/crawler-svc.pid"
PID=$(cat "$PID_FILE" 2>/dev/null || true)

if [ -n "$PID" ] && kill -0 "$PID" 2>/dev/null; then
  kill "$PID"
  sleep 1
  if kill -0 "$PID" 2>/dev/null; then
    kill -9 "$PID"
  fi
  echo "[ok] 已停止 crawler-svc (PID=$PID)"
else
  if pkill -f "bin/crawler-svc"; then
    echo "[ok] 已通过 pkill 清理 crawler-svc"
  else
    echo "[info] crawler-svc 未在运行"
  fi
fi

rm -f "$PID_FILE"