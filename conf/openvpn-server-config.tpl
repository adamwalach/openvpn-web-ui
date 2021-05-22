management {{ .Management }}

port {{ .Port }}
proto {{ .Proto }}

dev tun

ca {{ .Ca }}
cert {{ .Cert }}
key {{ .Key }}

cipher {{ .Cipher }}
keysize {{ .Keysize }}
auth {{ .Auth }}
dh {{ .Dh }}

server 10.9.0.0 255.255.255.0
ifconfig-pool-persist {{ .IfconfigPoolPersist }}
push "route 10.101.0.0 255.255.0.0 10.9.0.5 10"
push "route 10.102.0.0 255.255.0.0 10.9.0.5 20"
push "route 10.103.0.0 255.255.0.0 10.9.0.5 30"
push "route 10.31.0.0 255.255.0.0 10.9.0.5 40"
push "route 10.32.0.0 255.255.0.0 10.9.0.5 50"
push "route 10.33.0.0 255.255.0.0 10.9.0.5 60"
push "route 10.180.0.0 255.255.0.0 10.9.0.5 70"
push "route 10.203.0.0 255.255.0.0 10.9.0.5 80"
push "dhcp-option DNS 8.8.8.8"
push "dhcp-option DNS 8.8.4.4"

keepalive {{ .Keepalive }}

comp-lzo
max-clients {{ .MaxClients }}

persist-key
persist-tun

log         openvpn.log
verb 3

mute 10
