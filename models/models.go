package models

import (
	"os"

	"github.com/adamwalach/go-openvpn/server/config"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	passlib "gopkg.in/hlandau/passlib.v1"
)

var GlobalCfg Settings

func init() {
	initDB()
	createDefaultUsers()
	createDefaultSettings()
	createDefaultOVConfig()
}

func initDB() {
	orm.RegisterDriver("sqlite3", orm.DRSqlite)
	dbSource := "file:" + beego.AppConfig.String("dbPath")

	err := orm.RegisterDataBase("default", "sqlite3", dbSource)
	if err != nil {
		panic(err)
	}
	orm.Debug = true
	orm.RegisterModel(
		new(User),
		new(Settings),
		new(OVConfig),
	)

	// Database alias.
	name := "default"
	// Drop table and re-create.
	force := false
	// Print log.
	verbose := true

	err = orm.RunSyncdb(name, force, verbose)
	if err != nil {
		beego.Error(err)
		return
	}
}

func createDefaultUsers() {
	hash, err := passlib.Hash("b3secure")
	if err != nil {
		beego.Error("Unable to hash password", err)
	}
	user := User{
		Id:       1,
		Login:    "admin",
		Name:     "Administrator",
		Email:    "root@localhost",
		Password: hash,
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&user, "Name"); err == nil {
		if created {
			beego.Info("Default admin account created")
		} else {
			beego.Debug(user)
		}
	}

}

func createDefaultSettings() {
	s := Settings{
		Profile:       "default",
		MIAddress:     "openvpn:2080",
		MINetwork:     "tcp",
		ServerAddress: "127.0.0.1",
		OVConfigPath:  "/etc/openvpn/",
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&s, "Profile"); err == nil {
		GlobalCfg = s

		if created {
			beego.Info("New settings profile created")
		} else {
			beego.Debug(s)
		}
	} else {
		beego.Error(err)
	}
}

func createDefaultOVConfig() {
	c := OVConfig{
		Profile: "default",
		Config: config.Config{
			Port:                1194,
			Proto:               "udp",
			DNSServerOne:        "8.8.8.8",
			DNSServerTwo:        "8.8.4.4",
			Cipher:              "AES-256-CBC",
			Keysize:             256,
			Auth:                "SHA256",
			Dh:                  "dh2048.pem",
			Keepalive:           "10 120",
			IfconfigPoolPersist: "ipp.txt",
			Management:          "0.0.0.0 2080",
			MaxClients:          100,
			Server:              "10.8.0.0 255.255.255.0",
			Ca:                  "keys/ca.crt",
			Cert:                "keys/server.crt",
			Key:                 "keys/server.key",
			ExtraServerOptions:  "",
			ExtraClientOptions:  "",
		},
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&c, "Profile"); err == nil {
		if created {
			beego.Info("New settings profile created")
		} else {
			beego.Debug(c)
		}
		path := GlobalCfg.OVConfigPath + "/server.conf"
		if _, err = os.Stat(path); os.IsNotExist(err) {
			destPath := GlobalCfg.OVConfigPath + "/server.conf"
			if err = config.SaveToFile("conf/openvpn-server-config.tpl",
				c.Config, destPath); err != nil {
				beego.Error(err)
			}
		}
	} else {
		beego.Error(err)
	}
}
