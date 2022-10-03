#!/bin/sh

CURDIR=$(dirname "$0")
cd $CURDIR

PROJECTNODE="cronynode"

PROJECTBASE="."
ProjectBin=$PROJECTBASE/bin

NodeConf="$ProjectBin/node"
echo "start build file to $ProjectBin"

check() {
	EXCODE=$?
	if [ "$EXCODE" != "0" ]; then
		echo "build fail."
		exit $EXCODE
	fi
}
rm -rf $NodeConf


mkdir -p $NodeConf/logs
cp -r node/conf $NodeConf

#开启四个
echo "building project cronynode..."
go build -o $ProjectBin/$PROJECTNODE ./node/cmd/main.go
go build -o $ProjectBin/$PROJECTNODE"2" ./node/cmd/main.go
go build -o $ProjectBin/$PROJECTNODE"3" ./node/cmd/main.go
go build -o $ProjectBin/$PROJECTNODE"4" ./node/cmd/main.go
check



echo "build success."

