client
remote {{ .ServerAddress }} {{ .Port }}
proto {{ .Proto }}
dev tun
remote-cert-tls server
comp-lzo
;auth-user-pass
persist-key
persist-tun
nobind
resolv-retry infinite
verb 3
mute 10
<ca>
{{ .Ca }}
</ca>
<cert>
{{ .Cert }}
</cert>
<key>
{{ .Key }}
</key>