package models

import (
	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var engine *xorm.Engine

func init() {
	engine, _ = xorm.NewEngine("mysql", beego.AppConfig.String("explorer"))
	engine.SetMaxOpenConns(50)
}
