@echo off
setlocal

set C_PATH=%~dp0
for %%B in (%C_PATH%.) do set C_PATH=%%~dpB
for %%B in (%C_PATH%.) do set C_PATH=%%~dpB
echo GOPATH=%C_PATH%
set GOPATH=%C_PATH%

for /f "delims=" %%i in ('where bash') do set GITPATH=%%~dpi
echo GITPATH=%GITPATH%
set PATH=%GITPATH%;%PATH%
"%GITPATH%bash" ./test-wc.sh
REM "%GITPATH%sort.exe" -n -k2 mrtmp.wcseq
REM setlocal EnableDelayedExpansion
REM set Files=
REM for %%f in (pg-*.txt) do set Files=!Files! %%~f
REM go run wc.go master sequential %Files%