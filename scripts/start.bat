@echo off

set SERVICE_NAME=unitool_serve_windows.exe
if "%1" neq "" set SERVICE_NAME=%1
start /B %SERVICE_NAME% -port 8080 > output.log 2>&1