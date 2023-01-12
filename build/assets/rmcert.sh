#!/bin/bash
# Exit immediately if a command exits with a non-zero status
set -e

# .ovpn file path
DEST_FILE_PATH="/etc/openvpn/clients/$1.ovpn"

# Check if .ovpn file exists
if [[ ! -f $DEST_FILE_PATH ]]; then
    echo "User not found."
    exit 1
fi

# Remove user from OpenVPN
rm -f /etc/openvpn/pki/certs_by_serial/$2.pem
rm -f /etc/openvpn/pki/issued/$1.crt
rm -f /etc/openvpn/pki/private/$1.key
rm -f /etc/openvpn/pki/reqs/$1.req
rm -f /etc/openvpn/clients/$1.ovpn

# Fix index.txt by removing the user from the list following the serial number
sed -i'.bak' "/${2}/d" /usr/share/easy-rsa/pki/index.txt

echo 'Remove done!'
echo 'If you want to disconnect the user please restart the service using docker-compose restart openvpn.'
