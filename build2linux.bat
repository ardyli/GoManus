@echo off
echo 开始交叉编译...

:: 设置交叉编译环境变量
set CGO_ENABLED=0
set GOOS=linux
set GOARCH=amd64

:: 编译Go程序，假设源代码文件名为main.go，输出文件名为myapp
go build -o gomanus main.go

if %ERRORLEVEL% == 0 (
    echo ok，gomanus
) else (
    echo Err，请检查错误信息
)

:: 恢复默认Windows编译环境（可选）
set CGO_ENABLED=1
set GOOS=windows
set GOARCH=amd64

pause
