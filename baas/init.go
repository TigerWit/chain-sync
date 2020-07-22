package baas

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"

	"github.com/hyperledger/fabric-sdk-go/pkg/core/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
)

var (
	EngineB *xorm.Engine
	SDKInstance *fabsdk.FabricSDK
)

func init() {
	EngineB, _ := xorm.NewEngine("mysql", "root:@/fxchain")
	EngineB.SetMaxOpenConns(50)

	SDKInstance, _ = fabsdk.New(config.FromFile("./sdk.yaml"))
}
