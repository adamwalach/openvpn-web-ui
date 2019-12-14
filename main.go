package main

import (
	"os"
	"github.com/adamwalach/openvpn-web-ui/lib"
	_ "github.com/adamwalach/openvpn-web-ui/routers"
	"github.com/astaxie/beego"
)

func main() {
	lib.AddFuncMaps()
	beego.AppConfigPath = os.Getenv("CONFIG_FILE")
	beego.Run()
}
