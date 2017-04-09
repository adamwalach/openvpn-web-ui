package controllers

import (
	"errors"
	"html/template"
	"time"

	passlib "gopkg.in/hlandau/passlib.v1"

	"github.com/adamwalach/openvpn-web-ui/models"
	"github.com/astaxie/beego"
)

type LoginController struct {
	BaseController
}

func (c *LoginController) Login() {
	if c.IsLogin {
		c.Ctx.Redirect(302, c.URLFor("MainController.Get"))
		return
	}

	c.TplName = "login.html"
	c.Data["xsrfdata"] = template.HTML(c.XSRFFormHTML())
	if !c.Ctx.Input.IsPost() {
		return
	}

	flash := beego.NewFlash()
	login := c.GetString("login")
	password := c.GetString("password")

	user, err := Authenticate(login, password)
	if err != nil || user.Id < 1 {
		flash.Warning(err.Error())
		flash.Store(&c.Controller)
		return
	}
	flash.Success("Success logged in")
	flash.Store(&c.Controller)

	c.SetLogin(user)

	c.Redirect(c.URLFor("MainController.Get"), 303)
}

func (c *LoginController) Logout() {
	c.DelLogin()
	flash := beego.NewFlash()
	flash.Success("Success logged out")
	flash.Store(&c.Controller)

	c.Ctx.Redirect(302, c.URLFor("LoginController.Login"))
}

func Authenticate(login string, password string) (user *models.User, err error) {
	msg := "invalid login or password."
	user = &models.User{Login: login}

	if err := user.Read("Login"); err != nil {
		if err.Error() == "<QuerySeter> no row found" {
			err = errors.New(msg)
		}
		return user, err
	} else if user.Id < 1 {
		// No user
		return user, errors.New(msg)
		//} else if user.Password != password {
	} else if _, err := passlib.Verify(password, user.Password); err != nil {
		// No matched password
		return user, errors.New(msg)
	}
	user.Lastlogintime = time.Now()
	user.Update("Lastlogintime")
	return user, nil
}
