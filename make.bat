@echo off
title Auto Make

set FNAME=luzhibo
set PNAME=github.com\ggggle\%FNAME%
set GPATH=https://github.com/ggggle/luzhibo.git
set CPATH=%cd%
set BPATH=%~dp0



if "%1%"=="init" goto init

if "%GOPATH%"=="" set GOPATH=%UserProfile%\go
set Path=%Path%;%GOPATH%\bin
set SPATH=%GOPATH%\src\%PNAME%
git clone %GPATH% %SPATH%
cd %SPATH%
git pull
cd %CPATH%

if exist releases rd /s /q releases
md releases

::amd64:9
set GOARCH=amd64

set GOOS=linux
call:make
set GOOS=windows
call:make

:done
echo All done.
pause
goto:eof

:init
echo Initing...
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go get github.com/PuerkitoBio/goquery
go get github.com/pkg/browser
go get github.com/Baozisoftware/qrcode-terminal-go
go get github.com/lxn/walk
go get github.com/dkua/go-ico
go get github.com/inconshreveable/go-update
go get github.com/Baozisoftware/GoldenDaemon
go get github.com/Baozisoftware/golibraries
goto:done

:make
set TNAME=%FNAME%_%GOOS%_%GOARCH%
set LDFLAGS="-s -w"
if %GOOS%==windows set TNAME=%TNAME%.exe && go generate %PNAME% && set LDFLAGS="-s -w -H=windowsgui"
set TPATH=releases\%TNAME%
echo Building %TNAME%...
echo %LDFLAGS% ... %TPATH% ... %PNAME%
go build -ldflags=%LDFLAGS% -o %TPATH% %PNAME%
if exist %SPATH%\resource.syso del %SPATH%\resource.syso
%BPATH%upx --lzma --best -q %TPATH%
goto:eof
