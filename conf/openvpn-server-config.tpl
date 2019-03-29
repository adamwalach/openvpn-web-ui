management {{ .Management }}
verb 3

port {{ .Port }}
proto {{ .Proto }}

ca {{ .Ca }}
cert {{ .Cert }}
key {{ .Key }}

cipher {{ .Cipher }}
keysize {{ .Keysize }}
auth {{ .Auth }}
dh {{ .Dh }}

ifconfig-pool-persist {{ .IfconfigPoolPersist }}
server 192.168.255.0 255.255.255.0
### Route Configurations Below
route 192.168.254.0 255.255.255.0

### Push Configurations Below
push "block-outside-dns"
push "dhcp-option DNS 8.8.8.8"
push "dhcp-option DNS 8.8.4.4"
push "comp-lzo no"

dev tun
key-direction 0
keepalive {{ .Keepalive }}
persist-key
persist-tun
user nobody
group nogroup
comp-lzo no
mute 10

max-clients {{ .MaxClients }}

log  openvpn.log
