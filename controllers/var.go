package controllers

import (
	"github.com/astaxie/beego"
)

var (
	channelID = beego.AppConfig.String("channelid")
	user      = beego.AppConfig.String("user")
	org       = beego.AppConfig.String("org")
)
