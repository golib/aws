#!/usr/bin/env bash

if [ "$PKGROOT" == "" ]; then
    export PKGROOT=$(pwd)
fi

# adjust GOPATH
case ":$GOPATH:" in
    *":$PKGROOT:"*) :;;
    *) GOPATH=$PKGROOT:$GOPATH;;
esac
export GOPATH


# adjust PATH
if [ -n "$ZSH_VERSION" ]; then
    readopts="rA"
else
    readopts="ra"
fi
while IFS=':' read -$readopts ARR; do
    for i in "${ARR[@]}"; do
        case ":$PATH:" in
            *":$i/bin:"*) :;;
            *) PATH=$i/bin:$PATH
        esac
    done
done <<< "$GOPATH"
export PATH


# mock development && test envs
if [ ! -d "$PKGROOT/src/github.com/golib/aws" ]; then
    mkdir -p "$PKGROOT/src/github.com/golib"
    ln -s $PKGROOT "$PKGROOT/src/github.com/golib/aws"
fi
