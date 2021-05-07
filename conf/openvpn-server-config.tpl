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

server 10.8.0.0 255.255.255.0
ifconfig-pool-persist {{ .IfconfigPoolPersist }}
push "route 10.8.0.0 255.255.255.0"
push "dhcp-option DNS {{ .DNSServerOne }}"
push "dhcp-option DNS {{ .DNSServerTwo }}"

keepalive {{ .Keepalive }}

comp-lzo
max-clients {{ .MaxClients }}

persist-key
persist-tun

log         openvpn.log
verb 3

mute 10

{{ .ExtraServerOptions }}