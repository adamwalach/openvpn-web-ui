package controllers

import (
	"github.com/adamwalach/openvpn-web-ui/state"
	"html/template"
	"path/filepath"

	"github.com/adamwalach/go-openvpn/server/config"
	mi "github.com/adamwalach/go-openvpn/server/mi"
	"github.com/adamwalach/openvpn-web-ui/lib"
	"github.com/adamwalach/openvpn-web-ui/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type OVConfigController struct {
	BaseController
	ConfigDir string
}

func (c *OVConfigController) NestPrepare() {
	if !c.IsLogin {
		c.Ctx.Redirect(302, c.LoginPath())
		return
	}
	c.Data["breadcrumbs"] = &BreadCrumbs{
		Title: "OpenVPN config",
	}
}

func (c *OVConfigController) Get() {
	c.TplName = "ovconfig.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	cfg := models.OVConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	c.Data["Settings"] = &cfg

}

func (c *OVConfigController) Post() {
	c.TplName = "ovconfig.html"
	flash := beego.NewFlash()
	cfg := models.OVConfig{Profile: "default"}
	_ = cfg.Read("Profile")
	if err := c.ParseForm(&cfg); err != nil {
		beego.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}
	lib.Dump(cfg)
	c.Data["Settings"] = &cfg

	destPath := filepath.Join(state.GlobalCfg.OVConfigPath, "server.conf")
	err := config.SaveToFile(filepath.Join(c.ConfigDir, "openvpn-server-config.tpl"), cfg.Config, destPath)
	if err != nil {
		beego.Warning(err)
		flash.Error(err.Error())
		flash.Store(&c.Controller)
		return
	}

	o := orm.NewOrm()
	if _, err := o.Update(&cfg); err != nil {
		flash.Error(err.Error())
	} else {
		flash.Success("Config has been updated")
		client := mi.NewClient(state.GlobalCfg.MINetwork, state.GlobalCfg.MIAddress)
		if err := client.Signal("SIGTERM"); err != nil {
			flash.Warning("Config has been updated but OpenVPN server was NOT reloaded: " + err.Error())
		}
	}
	flash.Store(&c.Controller)
}
