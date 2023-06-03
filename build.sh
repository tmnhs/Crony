#!/bin/sh

CURDIR=$(dirname "$0")
cd $CURDIR

PROJECTADMIN="cronyadmin"
PROJECTNODE="cronynode"

PROJECTBASE="."
ProjectBin=$PROJECTBASE/bin

AdminConf="$ProjectBin/admin"
NodeConf="$ProjectBin/node"
WebFile="admin/web"
echo "start build file to $ProjectBin"

check() {
	EXCODE=$?
	if [ "$EXCODE" != "0" ]; then
		echo "build fail."
		exit $EXCODE
	fi
}
rm -rf $ProjectBin

mkdir -p $AdminConf/logs/
cp -r admin/conf $AdminConf

mkdir -p $NodeConf/logs
cp -r node/conf $NodeConf

#admin
echo "building project cronyadmin..."
go build -o $ProjectBin/$PROJECTADMIN ./admin/cmd/main.go
check

#node
echo "building project cronynode..."
go build -o $ProjectBin/$PROJECTNODE ./node/cmd/main.go
check

#web
echo "building web..."
cd $WebFile
npm install
yarn build
mv ./dist ../../bin/

echo "build success."

