// @Title
// @Description
// @Author  Wangwengang  2023/12/12 21:43
// @Update  Wangwengang  2023/12/12 21:43
package srpc

import (
	"fmt"
	"time"

	"github.com/rpcxio/rpcx-etcd/serverplugin"
	"github.com/smallnest/rpcx/protocol"
	"github.com/smallnest/rpcx/server"
	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/utils"
)

type SRPCServer interface {
	RegisterName(name string, rcvr interface{}, metadata string) error
	Serve(network, address string) (err error)
}

func AddRegistryPlugin(s *server.Server, rpc sconfig.RPC, service sconfig.RpcService) {

	r := &serverplugin.EtcdV3RegisterPlugin{
		ServiceAddress: fmt.Sprintf("tcp@%s:%s", service.ServiceAddr, service.Port),
		EtcdServers:    rpc.RegisterAddr,
		BasePath:       rpc.BasePath,
		UpdateInterval: time.Minute,
	}
	err := r.Start()
	if err != nil {
		panic(err)
	}
	protocol.Compressors[Brotli] = &utils.BrotliCompressor{}
	s.Plugins.Add(r)
}
