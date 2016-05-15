#!/bin/sh
SCRIPTPATH=$(cd "$(dirname "$0")"; pwd)

# set link

path="$SCRIPTPATH/src/github.com/leanote"
if [ ! -d "$path" ]; then
	mkdir -p "$path"
fi
rm -rf $SCRIPTPATH/src/github.com/lhboy1984/leanote # 先删除
ln -s ../../../../ $SCRIPTPATH/src/github.com/lhboy1984/leanote

# set GOPATH
export GOPATH=$SCRIPTPATH

script="$SCRIPTPATH/leanote-darwin-amd64"
chmod 777 $script
$script -importPath github.com/lhboy1984/leanote