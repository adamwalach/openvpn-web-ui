package main

import (
	"flag"
	"fmt"
	"github.com/adamwalach/openvpn-web-ui/lib"
	"github.com/adamwalach/openvpn-web-ui/models"
	"github.com/adamwalach/openvpn-web-ui/routers"
	_ "github.com/adamwalach/openvpn-web-ui/routers"
	"github.com/astaxie/beego"
	"path/filepath"
)

func main() {
	configDir := flag.String("config", "conf", "Path to config dir")
	flag.Parse()

	configFile := filepath.Join(*configDir, "app.conf")
	fmt.Println("Config file:", configFile)
	err := beego.LoadAppConfig("ini", configFile)

	models.InitDB()
	models.CreateDefaultUsers()
	models.CreateDefaultSettings()
	models.CreateDefaultOVConfig(*configDir)

	routers.Init()
	if err != nil {
		panic(err)
	}

	lib.AddFuncMaps()
	beego.Run()
}
