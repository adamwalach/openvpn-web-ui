package main

import (
	"flag"
	"fmt"
	"github.com/adamwalach/openvpn-web-ui/lib"
	_ "github.com/adamwalach/openvpn-web-ui/routers"
	"github.com/astaxie/beego"
)

func init() {
	fmt.Println("init")
	configFile := flag.String("config", "conf/app.conf", "Path to config file")
	flag.Parse()
	fmt.Println("Config file", *configFile)
	err := beego.LoadAppConfig("ini", *configFile)
	if err != nil {
		panic(err)
	}
}

func main() {
	lib.AddFuncMaps()
	beego.Run()
}
