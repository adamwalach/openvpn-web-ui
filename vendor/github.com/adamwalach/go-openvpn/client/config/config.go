package config

import (
	"bytes"
	"text/template"
	"io/ioutil"
)

var defaultConfig = Config{
	ServerAddress: "127.0.0.1",
	Port:          1194,
	Proto:         "udp",
	Cipher:        "AES-256-CBC",
	Keysize:       256,
	Auth:          "SHA256",
	Ca:            "ca.crt",
}

//Config model
type Config struct {
	ServerAddress string
	Port          int
	Proto         string

	Ca   string
	Cert string
	Key  string

	Cipher  string
	Keysize int
	Auth    string
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
