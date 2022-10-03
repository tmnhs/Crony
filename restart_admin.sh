#!/bin/sh

CURDIR=$(dirname "$0")
cd $CURDIR

PROJECTADMIN="cronyadmin"

PROJECTBASE="."
ProjectBin=$PROJECTBASE/bin
AdminConf="$ProjectBin/admin"

echo "start build file to $ProjectBin"

check() {
	EXCODE=$?
	if [ "$EXCODE" != "0" ]; then
		echo "build fail."
		exit $EXCODE
	fi
}
rm -rf $AdminConf

mkdir -p $AdminConf/logs/
cp -r admin/conf $AdminConf

echo "building project cronyadmin..."
go build -o $ProjectBin/$PROJECTADMIN ./admin/cmd/main.go
check

./server.sh restart admin