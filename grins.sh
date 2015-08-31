#!/bin/bash -
#

PROGRAM=$0
if test $# -eq 0
then
	echo "Usage: $0 <owner/repo>"
	exit 1
fi

SERVER="{{.}}"

if test ${#SERVER} -eq 5
then
    SERVER="grget.shengxiang.me"
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

getTarget(){
    set -eu
    URL="$SERVER/${1}/$REF/$OS/$ARCH" 
    TARGET=${1#*/}
    echo ">>> $URL"
    STATUS=$(curl --compressed -w %{http_code} "$URL" -o ${TARGET}.download)
    if test $STATUS -eq 200
    then
        chmod +x ${TARGET}.download
        mv ${TARGET}.download ${TARGET}
        echo "Save to ${TARGET}"
    else
        echo "EXIT: $?"
        cat ${TARGET}.download
        rm ${TARGET}.download
    fi
}

REPO=${1}
case "$REPO" in
    */*)
        getTarget "$REPO"
        ;;
    *)
        GUESSREPO=$(curl "$SERVER/lucky/${1}")
        echo $?
        if test -z "$GUESSREPO"
        then
            echo "Not found ${1}"
            exit 1
        fi
        getTarget "$GUESSREPO"
        ;;
esac
