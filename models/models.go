package models

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/adamwalach/go-openvpn/server/config"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"gopkg.in/hlandau/passlib.v1"
)

func InitDB() {
	err := orm.RegisterDriver("sqlite3", orm.DRSqlite)
	if err != nil {
		panic(err)
	}
	dbSource := "file:" + beego.AppConfig.String("dbPath")

	err = orm.RegisterDataBase("default", "sqlite3", dbSource)
	if err != nil {
		panic(err)
	}
	orm.Debug = true
	orm.RegisterModel(
		new(User),
		new(Settings),
		new(OVConfig),
	)

	err = orm.RunSyncdb("default", false, true)
	if err != nil {
		beego.Error(err)
		return
	}
}

func CreateDefaultUsers() {
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

func CreateDefaultSettings() (*Settings, error) {
	s := Settings{
		Profile:       "default",
		MIAddress:     beego.AppConfig.String("OpenVpnManagementAddress"),
		MINetwork:     beego.AppConfig.String("OpenVpnManagementNetwork"),
		ServerAddress: beego.AppConfig.String("OpenVpnServerAddress"),
		OVConfigPath:  beego.AppConfig.String("OpenVpnPath"),
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&s, "Profile"); err == nil {
		if created {
			beego.Info("New settings profile created")
		} else {
			beego.Debug(s)
		}
		return &s, nil
	} else {
		return nil, err
	}
}

func CreateDefaultOVConfig(configDir string, ovConfigPath string, address string, network string) {
	c := OVConfig{
		Profile: "default",
		Config: config.Config{
			Port:                1194,
			Proto:               "udp",
			Cipher:              "AES-256-CBC",
			Keysize:             256,
			Auth:                "SHA256",
			Dh:                  filepath.Join(ovConfigPath, "dh2048.pem"),
			Keepalive:           "10 120",
			IfconfigPoolPersist: "ipp.txt",
			Management:          fmt.Sprintf("%s %s", address, network),
			MaxClients:          100,
			Server:              "10.8.0.0 255.255.255.0",
			Ca:                  filepath.Join(ovConfigPath, "keys/ca.crt"),
			Cert:                filepath.Join(ovConfigPath, "keys/server.crt"),
			Key:                 filepath.Join(ovConfigPath, "keys/server.key"),
		},
	}
	o := orm.NewOrm()
	if created, _, err := o.ReadOrCreate(&c, "Profile"); err == nil {
		if created {
			beego.Info("New settings profile created")
		} else {
			beego.Debug(c)
		}
		serverConfig := filepath.Join(ovConfigPath, "server.conf")
		if _, err = os.Stat(serverConfig); os.IsNotExist(err) {
			if err = config.SaveToFile(filepath.Join(configDir, "openvpn-server-config.tpl"), c.Config, serverConfig); err != nil {
				beego.Error(err)
			}
		}
	} else {
		beego.Error(err)
	}
}
