@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo AI热点追踪平台 - Windows 停止脚本
echo ========================================
echo.

echo [1/3] 停止 Nginx...
taskkill /IM nginx.exe /F >nul 2>&1
if errorlevel 1 (
    echo [完成] Nginx 未运行或已停止
) else (
    echo [完成] Nginx 已停止
)
echo.

echo [2/3] 停止 Learning Path Service...
taskkill /IM learning-path-svc.exe /F >nul 2>&1
if errorlevel 1 (
    echo [完成] Learning Path Service 未运行或已停止
) else (
    echo [完成] Learning Path Service 已停止
)
echo.

echo [3/3] 停止 Content Service...
taskkill /IM content-svc.exe /F >nul 2>&1
if errorlevel 1 (
    echo [完成] Content Service 未运行或已停止
) else (
    echo [完成] Content Service 已停止
)
echo.

echo [4/4] 停止 Gateway...
taskkill /IM gateway.exe /F >nul 2>&1
if errorlevel 1 (
    echo [完成] Gateway 未运行或已停止
) else (
    echo [完成] Gateway 已停止
)
echo.

echo ========================================
echo 所有服务已停止！
echo ========================================
echo.
pause
