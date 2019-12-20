package main

import (
	"flag"
	"fmt"
	"github.com/adamwalach/openvpn-web-ui/lib"
	"github.com/adamwalach/openvpn-web-ui/models"
	"github.com/adamwalach/openvpn-web-ui/routers"
	_ "github.com/adamwalach/openvpn-web-ui/routers"
	"github.com/adamwalach/openvpn-web-ui/state"
	"github.com/astaxie/beego"
	"path/filepath"
)

func main() {
	configDir := flag.String("config", "conf", "Path to config dir")
	flag.Parse()

	configFile := filepath.Join(*configDir, "app.conf")
	fmt.Println("Config file:", configFile)

	if err := beego.LoadAppConfig("ini", configFile); err != nil {
		panic(err)
	}

	models.InitDB()
	models.CreateDefaultUsers()
	defaultSettings, err := models.CreateDefaultSettings()
	if err != nil {
		panic(err)
	}

	models.CreateDefaultOVConfig(*configDir, defaultSettings.OVConfigPath, defaultSettings.MIAddress, defaultSettings.MINetwork)

	state.GlobalCfg = *defaultSettings

	routers.Init(*configDir)

	lib.AddFuncMaps()
	beego.Run()
}
