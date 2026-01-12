// @Title
// @Description
// @Author  Wangwengang  2023/12/10 11:46
// @Update  Wangwengang  2023/12/10 11:46
package server

import (
	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/snet/http"
	"github.com/wwengg/threego/core/srpc"
)

type GatewayOptions func(g *Gateway)

type Gateway struct {
	config     *sconfig.Gateway
	httpServer http.HttpServer

	RpcMode bool
	RPC     srpc.SRPC
}

func WithSRPC(srpc srpc.SRPC) GatewayOptions {
	return func(g *Gateway) {
		g.setRpc(srpc)
	}
}

func NewGateway(config *sconfig.Gateway, httpServer http.HttpServer, opts ...GatewayOptions) *Gateway {
	g := &Gateway{
		config:     config,
		httpServer: httpServer,
		RpcMode:    false,
	}

	for _, opt := range opts {
		opt(g)
	}
	return g
}

func (g *Gateway) Start() {
	g.httpServer.Serve()

}

func (g *Gateway) setRpc(srpc srpc.SRPC) {
	g.RpcMode = true
	g.RPC = srpc
}
