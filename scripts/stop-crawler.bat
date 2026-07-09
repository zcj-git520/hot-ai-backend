@echo off
chcp 65001 >nul

tasklist /FI "IMAGENAME eq crawler-svc.exe" 2>nul | find /I "crawler-svc.exe" >nul
if errorlevel 1 (
  echo [info] crawler-svc 未在运行
  exit /b 0
)

taskkill /F /IM crawler-svc.exe
echo [ok] 已停止 crawler-svc