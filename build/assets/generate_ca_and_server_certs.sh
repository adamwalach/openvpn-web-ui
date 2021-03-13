#!/bin/bash -ex

CA_NAME=LocalCA
EASY_RSA=/usr/share/easy-rsa
OD=$PWD
export EASYRSA_BATCH="true"

dd if=/dev/urandom of=/etc/openvpn/pki/.rnd bs=256 count=1

cd $OVDIR

export KEY_NAME=$CA_NAME
echo "Generating CA cert"
$EASY_RSA/easyrsa init-pki
cp -f /opt/scripts/vars.template $OVDIR/pki/vars 
dd if=/dev/urandom of=/etc/openvpn/pki/.rnd bs=256 count=1 > /dev/null 2>&1

$EASY_RSA/easyrsa build-ca nopass

# only temporarily for tests as it takes ages... use existing one

# $EASY_RSA/easyrsa gen-dh

#$EASY_RSA/build-ca
#export EASY_RSA="${EASY_RSA:-.}"

# build server key
echo "Generating server cert $SERVER_FQDN"
export KEY_NAME=$SERVER_FQDN
$EASY_RSA/easyrsa build-server-full $SERVER_FQDN nopass

$EASY_RSA/easyrsa gen-crl

echo "Missing is still ta.key"
echo "openvpn --genkey --secret /root/easy-rsa-example/pki/ta.key"
