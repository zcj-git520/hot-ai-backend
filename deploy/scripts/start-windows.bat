@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo ========================================
echo AI热点追踪平台 - Windows 启动脚本
echo ========================================
echo.

REM 检查 Go 是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未检测到 Go，请先安装 Go 1.22+
    pause
    exit /b 1
)



REM 检查配置文件
if not exist "config\app.env" (
    echo [提示] 未找到配置文件，正在创建默认配置...
    if not exist "config" mkdir config
    (
        echo DB_DSN=root:zcj@tcp123456@tcp(localhost:3306)/hot_ai?charset=utf8mb4&parseTime=True&loc=Local
        echo CONTENT_SVC_URL=http://localhost:8001
        echo JWT_SECRET=your-jwt-secret-key-change-in-production
        echo REDIS_URL=localhost:6379
    ) > "config\app.env"
    echo [完成] 默认配置文件已创建
    echo [警告] 请修改 config\app.env 中的数据库密码和 JWT Secret
    echo.
)

REM 创建 logs 目录
if not exist "logs" mkdir logs

echo [1/6] 编译 Gateway 服务...
cd /d "%~dp0..\.."
go build -o bin\gateway.exe ./apps/gateway || exit /b 1
echo [完成] Gateway 编译成功
echo.

echo [2/6] 编译 Content Service...
go build -o bin\content-svc.exe ./apps/services/content-svc || exit /b 1
echo [完成] Content Service 编译成功
echo.

echo [3/6] 编译 Learning Path Service...
go build -o bin\learning-path-svc.exe ./apps/services/learning-path-svc || exit /b 1
echo [完成] Learning Path Service 编译成功
echo.

echo [4/6] 启动 Gateway (端口 8000)...
start "Gateway" cmd /c "bin\gateway.exe -f apps\gateway\etc\gateway.yaml > logs\gateway.log 2>&1"
timeout /t 3 /nobreak >nul
echo [完成] Gateway 已启动
echo.

echo [5/6] 启动 Content Service (端口 8001)...
start "Content-SVC" cmd /c "bin\content-svc.exe -f apps\services\content-svc\etc\content-svc.yaml > logs\content-svc.log 2>&1"
timeout /t 3 /nobreak >nul
echo [完成] Content Service 已启动
echo.

echo [6/6] 启动 Learning Path Service (端口 8003)...
start "Learning-Path-SVC" cmd /c "bin\learning-path-svc.exe -f apps\services\learning-path-svc\etc\learning-path-svc.yaml > logs\learning-path-svc.log 2>&1"
timeout /t 3 /nobreak >nul
echo [完成] Learning Path Service 已启动
echo.


echo ========================================
echo 所有服务启动完成！
echo ========================================
echo.
echo 服务地址：
echo   - Gateway API: http://localhost:8000
echo   - Content API: http://localhost:8001
echo   - Learning Path API: http://localhost:8003
if exist "nginx.exe" (
    echo   - Nginx Proxy: http://localhost
)
echo.
echo 日志文件位置：
echo   - Gateway: logs\gateway.log
echo   - Content Service: logs\content-svc.log
echo   - Learning Path Service: logs\learning-path-svc.log
if exist "nginx.exe" (
    echo   - Nginx: logs\nginx\access.log, logs\nginx\error.log
)
echo.
echo 查看日志：
echo   - Gateway: type logs\gateway.log
echo   - Content Service: type logs\content-svc.log
echo.
echo 停止服务：
echo   - 关闭所有命令行窗口
echo   - 或运行 stop-windows.bat
echo.
pause
