package config

import (
	"bytes"
	"html/template"
	"io/ioutil"
)

var defaultConfig = Config{
	Port:                1194,
	Proto:               "udp",
	Cipher:              "AES-256-CBC",
	Keysize:             256,
	Auth:                "SHA256",
	Dh:                  "dh2048.pem",
	Keepalive:           "10 120",
	IfconfigPoolPersist: "ipp.txt",
}

//Config model
type Config struct {
	Port  int
	Proto string

	Ca   string
	Cert string
	Key  string

	Cipher  string
	Keysize int
	Auth    string
	Dh      string

	Server              string
	IfconfigPoolPersist string
	Keepalive           string
	MaxClients          int

	Management string
}

//New returns config object with default values
func New() Config {
	return defaultConfig
}

//GetText injects config values into template
func GetText(tpl string, c Config) (string, error) {
	t := template.New("config")
	t, err := t.Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	t.Execute(buf, c)
	return buf.String(), nil
}

//SaveToFile reads teamplate and writes result to destination file
func SaveToFile(tplPath string, c Config, destPath string) error {
	template, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return err
	}

	str, err := GetText(string(template), c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destPath, []byte(str), 0644)
}
