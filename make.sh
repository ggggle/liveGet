#!/bin/sh

echo -n "\033]0;Auto Make\007"

FNAME=liveGet
PNAME=github.com/ggggle/$FNAME
GPATH=https://github.com/ggggle/liveGet.git
CPATH=`pwd`
BPATH=`dirname $0`
UPX=$BPATH/upx
chmod +x $UPX

MAKE()
{
	TNAME="$FNAME"_"$GOOS"_"$GOARCH"
	LDFLAGS="-s -w"
	if [ "$GOOS" = "windows" ]; then
		TNAME=$TNAME.exe
		LDFLAGS="-s -w -H=windowsgui"
		GOOS=$GOOS GOARCH=$GOARCH go generate $PNAME
	fi
	TPATH=releases/$TNAME
	echo Building $TNAME...
    GOOS=$GOOS GOARCH=$GOARCH go build -ldflags="$LDFLAGS" -o $TPATH $PNAME
	if [ -f "$SPATH/resource.syso" ]; then
        rm $SPATH/resource.syso
    fi
    $UPX --lzma --best -q $TPATH
}

DONE()
{
	echo All done.
	exit 0
}

if [ "$GOPATH" = "" ]; then 
	GOPATH=~/go
fi

mkdir -p $GOPATH/src/github.com/Baozisoftware
cd $GOPATH/src/github.com/Baozisoftware
git clone https://github.com/Baozisoftware/golibraries
cd $GOPATH/src/github.com/Baozisoftware/golibraries
git checkout 91a9f7051cb37b11b3bd7bd16ffe0875e0e7de2e
cd $GOPATH/src/github.com/Baozisoftware
git clone https://github.com/Baozisoftware/qrcode-terminal-go
git clone https://github.com/Baozisoftware/GoldenDaemon

#init
echo Initing...
go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo
go get github.com/PuerkitoBio/goquery
go get github.com/pkg/browser
go get github.com/mattn/go-isatty
go get github.com/lxn/walk
go get github.com/dkua/go-ico
go get gopkg.in/Knetic/govaluate.v3
go get github.com/lxn/win
go get github.com/inconshreveable/go-update
go get github.com/buger/jsonparser
if [ "$1" = "init" ]; then
	DONE
fi

PATH=$PATH:$GOPATH/bin
SPATH=$GOPATH/src/$PNAME
git clone $GPATH $SPATH
cd $SPATH
git pull
cd $CPATH


if [ -d releases ]; then
	rm -rf releases
fi
mkdir releases

#386:7
GOARCH=386
GOOS=linux
MAKE
GOOS=windows
MAKE

#amd64:9
GOARCH=amd64
GOOS=linux
MAKE
GOOS=windows
MAKE

DONE

