package controllers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/adamwalach/go-openvpn/client/config"
	"github.com/adamwalach/openvpn-web-ui/lib"
	"github.com/adamwalach/openvpn-web-ui/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
)

type NewCertParams struct {
	Name string `form:"Name" valid:"Required;"`
}

type CertificatesController struct {
	BaseController
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
	filename := fmt.Sprintf("%s.zip", name)

	c.Ctx.Output.Header("Content-Type", "application/zip")
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	zw := zip.NewWriter(c.Controller.Ctx.ResponseWriter)

	keysPath := models.GlobalCfg.OVConfigPath + "keys/"
	if cfgPath, err := saveClientConfig(name); err == nil {
		addFileToZip(zw, cfgPath)
	}
	if ovpnPath, err := saveClientOvpn(name); err == nil {
		addFileToZip(zw, ovpnPath)
	}

	addFileToZip(zw, keysPath+"ca.crt")
	addFileToZip(zw, keysPath+name+".crt")
	addFileToZip(zw, keysPath+name+".key")

	if err := zw.Close(); err != nil {
		beego.Error(err)
	}
}

func addFileToZip(zw *zip.Writer, path string) error {
	header := &zip.FileHeader{
		Name:         filepath.Base(path),
		Method:       zip.Store,
		ModifiedTime: uint16(time.Now().UnixNano()),
		ModifiedDate: uint16(time.Now().UnixNano()),
	}
	fi, err := os.Open(path)
	if err != nil {
		beego.Error(err)
		return err
	}

	fw, err := zw.CreateHeader(header)
	if err != nil {
		beego.Error(err)
		return err
	}

	if _, err = io.Copy(fw, fi); err != nil {
		beego.Error(err)
		return err
	}

	return fi.Close()
}

// @router /certificates [get]
func (c *CertificatesController) Get() {
	c.TplName = "certificates.html"
	c.showCerts()
}

func (c *CertificatesController) showCerts() {
	path := models.GlobalCfg.OVConfigPath + "keys/index.txt"
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
			if err := lib.CreateCertificate(cParams.Name); err != nil {
				beego.Error(err)
				flash.Error(err.Error())
				flash.Store(&c.Controller)
			}
		}
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

func saveClientConfig(name string) (string, error) {
	cfg := config.New()
	cfg.ServerAddress = models.GlobalCfg.ServerAddress
	cfg.Cert = name + ".crt"
	cfg.Key = name + ".key"
	serverConfig := models.OVConfig{Profile: "default"}
	serverConfig.Read("Profile")
	cfg.Port = serverConfig.Port
	cfg.Proto = serverConfig.Proto
	cfg.Auth = serverConfig.Auth
	cfg.Cipher = serverConfig.Cipher
	cfg.Keysize = serverConfig.Keysize

	destPath := models.GlobalCfg.OVConfigPath + "keys/" + name + ".conf"
	if err := config.SaveToFile("conf/openvpn-client-config.tpl",
		cfg, destPath); err != nil {
		beego.Error(err)
		return "", err
	}

	return destPath, nil
}

func saveClientOvpn(name string) (string, error) {
	cfg := config.New()
	cfg.ServerAddress = models.GlobalCfg.ServerAddress
	serverConfig := models.OVConfig{Profile: "default"}
	serverConfig.Read("Profile")
	cfg.Port = serverConfig.Port
	cfg.Proto = serverConfig.Proto
	cfg.Auth = serverConfig.Auth
	cfg.Cipher = serverConfig.Cipher
	cfg.Keysize = serverConfig.Keysize

	keysPath := models.GlobalCfg.OVConfigPath + "keys/"
	caFilePath := keysPath + "ca.crt"
	certFilePath := keysPath + name + ".crt"
	keyFilePath := keysPath + name + ".key"

	if caByte, err := ioutil.ReadFile(caFilePath); err == nil {
		cfg.Ca = string(caByte)
	}
	if certByte, err := ioutil.ReadFile(certFilePath); err == nil {
		cfg.Cert = string(certByte)
	}
	if keyByte, err := ioutil.ReadFile(keyFilePath); err == nil {
		cfg.Key = string(keyByte)
	}

	destPath := models.GlobalCfg.OVConfigPath + "keys/" + name + ".ovpn"
	if err := saveToFile("conf/openvpn-client-ovpn.tpl",
		cfg, destPath); err != nil {
		beego.Error(err)
		return "", err
	}

	return destPath, nil
}

//SaveToFile reads teamplate and writes result to destination file  with text/template
func saveToFile(tplPath string, c config.Config, destPath string) error {
	templateByte, err := ioutil.ReadFile(tplPath)
	if err != nil {
		return err
	}

	t := template.New("config")
	temp, err := t.Parse(string(templateByte))
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	temp.Execute(buf, c)

	str := buf.String()
	fmt.Printf(str)
	return ioutil.WriteFile(destPath, []byte(str), 0644)
}
