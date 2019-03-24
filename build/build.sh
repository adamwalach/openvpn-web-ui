#!/bin/bash

# pack project by beego
# cd ../ && bee pack -be GOOS=linux

set -e

PKGFILE=openvpn-web-ui.tar.gz

cp -f ../$PKGFILE ./

docker build -t awalach/openvpn-web-ui .

rm -f $PKGFILE
