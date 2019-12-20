dev tun
persist-tun
persist-key
client
resolv-retry infinite
remote {{ .ServerAddress }} {{ .Port }} {{ .Proto }}
lport 0
cipher {{ .Cipher }}
auth {{ .Auth }}
tls-client
redirect-gateway def1
comp-lzo
<ca>
{{ .Ca }}
</ca>
<cert>
{{ .Cert }}
</cert>
<key>
{{ .Key }}
</key>