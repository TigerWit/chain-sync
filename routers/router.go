package routers

import (
	"chain-sync/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/getchannelinfo", &controllers.MainController{}, "get:GetChannelInfo")
	beego.Router("/getlastblocks", &controllers.MainController{}, "get:GetLastBlocks")
	beego.Router("/getlasttxs", &controllers.MainController{}, "get:GetLastTxs")
	beego.Router("/gettxsbyblocknum", &controllers.MainController{}, "get:GetTxsByBlocknum")
	beego.Router("/getoriginbyhash", &controllers.MainController{}, "get:GetOriginByHash")
}
