@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

cd /d "%~dp0\.."

if not exist logs mkdir logs
if not exist bin (
  echo [build] go build crawler-svc ...
  go build -o bin\crawler-svc.exe apps\crawler-svc
  if errorlevel 1 (
    echo [fail] 编译失败
    exit /b 1
  )
)

tasklist /FI "IMAGENAME eq crawler-svc.exe" 2>nul | find /I "crawler-svc.exe" >nul
if not errorlevel 1 (
  echo [skip] crawler-svc 已在运行
  exit /b 0
)

start /B "" cmd /c "bin\crawler-svc.exe -f apps\crawler-svc\etc\crawler-svc.yaml > logs\crawler.log 2>&1"

timeout /t 2 /nobreak >nul
tasklist /FI "IMAGENAME eq crawler-svc.exe" 2>nul | find /I "crawler-svc.exe" >nul
if not errorlevel 1 (
  echo [ok] crawler-svc 已后台启动，日志: logs\crawler.log
) else (
  echo [fail] crawler-svc 启动失败，查看 logs\crawler.log
  exit /b 1
)