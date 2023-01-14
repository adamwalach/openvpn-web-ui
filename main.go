package main

import (
	"flag"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/d3vilh/openvpn-web-ui/lib"
	"github.com/d3vilh/openvpn-web-ui/models"
	"github.com/d3vilh/openvpn-web-ui/routers"
	_ "github.com/d3vilh/openvpn-web-ui/routers"
	"github.com/d3vilh/openvpn-web-ui/state"
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
