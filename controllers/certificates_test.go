package controllers

import (
	"github.com/adamwalach/go-openvpn/client/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNoExcape(t *testing.T) {

	cfg := config.New()
	cfg.Cert = "ce++rt"

	text, err := GetText("++ {{ .Cert }} ++", cfg)
	assert.Equal(t, text, "++ ce++rt ++", err)

}
