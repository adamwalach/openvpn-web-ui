package main

import (
	"flag"
	"fmt"
	"github.com/adamwalach/openvpn-web-ui/lib"
	"github.com/adamwalach/openvpn-web-ui/models"
	_ "github.com/adamwalach/openvpn-web-ui/routers"
	"github.com/astaxie/beego"
	"path/filepath"
)

func init() {
	fmt.Println("init")
	configFile := flag.String("config", "conf/app.conf", "Path to config file")
	flag.Parse()
	configDir := filepath.Dir(*configFile)
	fmt.Println("Config file", *configFile)
	err := beego.LoadAppConfig("ini", *configFile)
	models.Init(configDir)
	if err != nil {
		panic(err)
	}
}

func main() {
	lib.AddFuncMaps()
	beego.Run()
}
