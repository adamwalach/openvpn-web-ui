package controllers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"text/template"

	"github.com/adamwalach/go-openvpn/client/config"
	"github.com/adamwalach/openvpn-web-ui/lib"
	"github.com/adamwalach/openvpn-web-ui/models"
	"github.com/adamwalach/openvpn-web-ui/state"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type NewCertParams struct {
	Name       string `form:"Name" valid:"Required;"`
	Staticip   string `form:"staticip"`
	Passphrase string `form:"passphrase"`
}

type CertificatesController struct {
	BaseController
	ConfigDir string
}

func (c *CertificatesController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "Certificates",
	}
}

// @router /certificates/:key [get]
func (c *CertificatesController) Download() {
	name := c.GetString(":key")
	filename := fmt.Sprintf("%s.ovpn", name)

	c.Ctx.Output.Header("Content-Type", "application/octet-stream")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	keysPath := filepath.Join(state.GlobalCfg.OVConfigPath, "pki/issued")

	cfgPath, err := c.saveClientConfig(keysPath, name)
	if err != nil {
		beego.Error(err)
		return
	}
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		beego.Error(err)
		return
	}
	if _, err = c.Controller.Ctx.ResponseWriter.Write(data); err != nil {
		beego.Error(err)
	}
}

// @router /certificates [get]
func (c *CertificatesController) Get() {
	c.TplName = "certificates.html"
	c.showCerts()
}

func (c *CertificatesController) showCerts() {
	path := filepath.Join(state.GlobalCfg.OVConfigPath, "pki/index.txt")
	certs, err := lib.ReadCerts(path)
	if err != nil {
		beego.Error(err)
	}
	lib.Dump(certs)
	c.Data["certificates"] = &certs
}

// @router /certificates [post]
func (c *CertificatesController) Post() {
	c.TplName = "certificates.html"
	flash := beego.NewFlash()

	cParams := NewCertParams{}
	if err := c.ParseForm(&cParams); err != nil {
		beego.Error(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
	} else {
		if vMap := validateCertParams(cParams); vMap != nil {
			c.Data["validation"] = vMap
		} else {
			if err := lib.CreateCertificate(cParams.Name, cParams.Staticip, cParams.Passphrase); err != nil {
				beego.Error(err)
				flash.Error(err.Error())
				flash.Store(&c.Controller)
			} else {
				flash.Success("Success! Certificate for the name \"" + cParams.Name + "\" has been created")
				flash.Store(&c.Controller)
			}
		}
	}
	c.showCerts()
}

// @router /certificates/revoke/:key [get]
func (c *CertificatesController) Revoke() {
	c.TplName = "certificates.html"
	flash := beego.NewFlash()
	name := c.GetString(":key")
	if err := lib.RevokeCertificate(name); err != nil {
		beego.Error(err)
		//flash.Error(err.Error())
		//flash.Store(&c.Controller)
	} else {
		flash.Warning("Success! Certificate for the name \"" + name + "\" has been revoked")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

// @router /certificates/restart [get]
func (c *CertificatesController) Restart() {
	lib.Restart()
	c.Redirect(c.URLFor("CertificatesController.Get"), 302)
	return
}

// @router /certificates/burn/:key/:serial [get]
func (c *CertificatesController) Burn() {
	c.TplName = "certificates.html"
	flash := beego.NewFlash()
	CN := c.GetString(":key")
	serial := c.GetString(":serial")
	if err := lib.BurnCertificate(CN, serial); err != nil {
		beego.Error(err)
		//flash.Error(err.Error())
		//flash.Store(&c.Controller)
	} else {
		flash.Success("Success! Certificate for the name \"" + CN + "\" has been removed")
		flash.Store(&c.Controller)
	}
	c.showCerts()
}

func validateCertParams(cert NewCertParams) map[string]map[string]string {
	valid := validation.Validation{}
	b, err := valid.Valid(&cert)
	if err != nil {
		beego.Error(err)
		return nil
	}
	if !b {
		return lib.CreateValidationMap(valid)
	}
	return nil
}

func (c *CertificatesController) saveClientConfig(keysPath string, name string) (string, error) {
	cfg := config.New()
	keysPathCa := filepath.Join(state.GlobalCfg.OVConfigPath, "pki")
	cfg.ServerAddress = state.GlobalCfg.ServerAddress
	ca, err := ioutil.ReadFile(filepath.Join(keysPathCa, "ca.crt"))
	if err != nil {
		return "", err
	}
	cfg.Ca = string(ca)
	cert, err := ioutil.ReadFile(filepath.Join(keysPath, name+".crt"))
	if err != nil {
		return "", err
	}
	cfg.Cert = string(cert)
	keysPathKey := filepath.Join(state.GlobalCfg.OVConfigPath, "pki/private")
	key, err := ioutil.ReadFile(filepath.Join(keysPathKey, name+".key"))
	if err != nil {
		return "", err
	}
	cfg.Key = string(key)
	serverConfig := models.OVConfig{Profile: "default"}
	_ = serverConfig.Read("Profile")
	cfg.Port = serverConfig.Port
	cfg.Proto = serverConfig.Proto
	cfg.Auth = serverConfig.Auth
	cfg.Cipher = serverConfig.Cipher
	cfg.Keysize = serverConfig.Keysize

	destPath := filepath.Join(state.GlobalCfg.OVConfigPath, "clients", name+".ovpn")
	if err := SaveToFile(filepath.Join(c.ConfigDir, "openvpn-client-config.tpl"), cfg, destPath); err != nil {
		beego.Error(err)
		return "", err
	}

	return destPath, nil
}

func GetText(tpl string, c config.Config) (string, error) {
	t := template.New("config")
	t, err := t.Parse(tpl)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	err = t.Execute(buf, c)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func SaveToFile(tplPath string, c config.Config, destPath string) error {
	tpl, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return err
	}

	str, err := GetText(string(tpl), c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(destPath, []byte(str), 0644)
}
