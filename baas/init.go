package baas

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var EngineB *xorm.Engine

func init() {
	EngineB, _ := xorm.NewEngine("mysql", "root:@/fxchain")
	EngineB.SetMaxOpenConns(50)
}
