# ThreeGO
敏捷开发、高效、稳定、易部署、多语言(计划Java、Python)

## 快速开始

### simplectl 快速生成rpc服务代码
1. install simplectl
```
go install github.com/wwengg/threego/tool/threegoctl@latest
```

2. Create a new directory
3. `cd` into that directory
4. run `go mod init <MODNAME>`

e.g.
```
mkdir myapp 
cd myapp
go mod init github.com/wwengg/myapp
simplectl rpc init --author "wwengg info@wwengg.cn"
go run main.go
```

### simplectl 快速生成基础增删改查proto文件
e.g
```
simplectl proto new user
```

### protoc 快速生成model和service
1. install simplectl
   `go install github.com/wwengg/protoc-gen-simple@latest`

e.g
```
protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative --simple_out=model --simple_opt=paths=source_relative pbuser/pbuser.proto
```

## TODO List

### gateway
- [x] http(example [im](https://github.com/wwengg/im))
- [x] tcp(example [im](https://github.com/wwengg/im))
- [x] kcp(example [im](https://github.com/wwengg/im))
- [x] websocket(example [im](https://github.com/wwengg/im))
- [x] quic
- [x] 智能路由 (example [im](https://github.com/wwengg/im))
- [ ] simplectl 快速生成gateway代码
- [ ] 限流
- [ ] 自动熔断
- [x] 链路追踪
- [x] 鉴权
- [ ] 验签加解密
- [x] 超时控制
- [ ] 监控报警
- [ ] 支持Json/Protobuf数据解析
 
### rpc service
- [x] simplectl 初始化项目(Tanks [Cobra](https://https://github.com/spf13/cobra))
- [x] simplectl 快速生成proto
- [x] simplectl 根据proto快速生成model([gorm](https://github.com/go-gorm/gorm))、service([rpcx-server](https://github.com/smallnest/rpcx))
- [x] 服务性能监控报警
- [x] 日志记录
- [ ] 调用鉴权
- [ ] 链路追踪
- [ ] 自动熔断
- [x] Java(Thanks[java-rpcx](https://github.com/smallnest/rpcx-java))
- [ ] Python

### middleware
- [x] etcd 分布式锁
- [x] etcd 服务发现
- [ ] k8s 服务发现
- [ ] Nsq消息队列
- [x] metrics

### Store
- [x] Mysql
- [ ] Redis
- [ ] MongoDB

## Thanks
- [rpcx](https://github.com/smallnest/rpcx)
- [cobra](https://https://github.com/spf13/cobra)