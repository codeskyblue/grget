#!/bin/bash -
#

PROGRAM=$0
if test $# -eq 0
then
	echo "Usage: $0 <owner/repo>"
	exit 1
fi

OS=$(uname | tr A-Z a-z)
ARCH=
case "$(uname -m)" in
	x86_64)
		ARCH=amd64
		;;
	i686)
		ARCH=i386
		;;
	*)
		echo "Usage: Architecture not found by uname -m"
		exit 2
		;;
esac
REF=${REF:-"master"}
TARGET=${1#*/}

set -eu
URL="grget.shengxiang.me/${1}/$REF/$OS/$ARCH" 
echo ">>> $URL"
curl "$URL" -o $TARGET
chmod +x $TARGET
