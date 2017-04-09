package main

import (
	"github.com/adamwalach/openvpn-web-ui/lib"
	_ "github.com/adamwalach/openvpn-web-ui/routers"
	"github.com/astaxie/beego"
)

func main() {
	lib.AddFuncMaps()
	beego.Run()
}
