#!/bin/bash -x

set -e
export OVDIR=/etc/openvpn

cd /opt/

if [ ! -f $OVDIR/.provisioned ]; then
  echo "Preparing certificates"
  mkdir -p $OVDIR/pki
  ./scripts/generate_ca_and_server_certs.sh
#  openssl dhparam -dsaparam -out $OVDIR/dh2048.pem 2048
  openssl dhparam -dsaparam -out $OVDIR/dh2048.pem 2048
  openssl dhparam -dsaparam -out $OVDIR/dh4096.pem 4096
#  touch $OVDIR/dh4096.pem
  cd $OVDIR/pki/
  ln -s ../dh4096.pem dh.pem
  touch $OVDIR/.provisioned
fi
cd /opt/openvpn-gui
mkdir -p db
./openvpn-web-ui
echo "Starting!"

