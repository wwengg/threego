// code generator by simplectl
package main

import (
	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/slog"
	"github.com/wwengg/threego/core/snet/http"
	"github.com/wwengg/threego/core/srpc"
	"github.com/wwengg/threego/server"
	"github.com/wwengg/threego/server/example/gateway/global"
	"github.com/wwengg/threego/server/example/gateway/middleware"
	"github.com/wwengg/threego/server/example/gateway/router"
)



func main() {
	// 初始化配置文件
	sconfig.Viper("./server/example/gateway/config.yaml")

	// 初始化日志
	global.SLog = slog.NewZapLog(&sconfig.S_CONF.Slog)

	// 初始化SRPC
	global.SRPC = srpc.NewSRPCClients(&sconfig.S_CONF.RPC)

	// 初始化gin
	ginEngine := http.NewGinEngine(&sconfig.S_CONF)

	// 配置路由，中间件
	publicGroup := ginEngine.GetPublicRouterGroup()
	privateGroup := ginEngine.GetPrivateRouterGroup()
	publicGroup.Use(middleware.BaseHandler())
	{
		router.InitSRPCRouter(publicGroup)
	}
	privateGroup.Use(middleware.BaseHandler())
	{
		router.InitSRPCRouter(privateGroup)
	}

	srv := server.NewGateway(&sconfig.S_CONF.Gateway, ginEngine)

	srv.Start()
}