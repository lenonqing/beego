package main

import (
	"framework/graylog"
	_ "framework/routers"

	"framework/db"

	"github.com/astaxie/beego"
)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	graylog.Init()
	db.Init()
	beego.Run()
}
