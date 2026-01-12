# Threego

[English](README.md) | [简体中文](README_CN.md)

一个高性能、生产级的 Go 微服务框架。Threego 提供了构建可扩展、易维护的分布式系统所需的完整工具链。

## 特性

- **多协议支持**：HTTP、WebSocket、QUIC、KCP、TCP
- **服务网格**：内置服务发现、负载均衡和熔断机制
- **高性能**：基于 rpcx 框架实现高效 RPC 通信
- **开发友好**：提供 CLI 工具 `simplectl` 快速生成项目脚手架
- **可观测性**：集成 Jaeger 链路追踪、结构化日志和 Sentry 错误监控
- **热更新**：使用 tableflip 实现零宕机部署
- **存储无关**：开箱即用支持 MySQL、Redis、MongoDB

## 架构

```
┌─────────────────────────────────────────────────────────────┐
│                      网关层                                  │
│              (HTTP/WebSocket/QUIC/KCP)                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    服务总线                                  │
│            (消息路由与分发)                                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    RPC 服务                                  │
│              (业务逻辑处理)                                  │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   存储层                                    │
│           (MySQL / Redis / MongoDB)                         │
└─────────────────────────────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│              服务发现 (etcd)                                │
│              配置中心                                        │
└─────────────────────────────────────────────────────────────┘
```

## 项目结构

```
threego/
├── core/                    # 核心框架模块
│   ├── internal/            # 内部常量和类型定义
│   ├── plugin/              # 插件系统 (Jaeger 等)
│   ├── sbus/                # 服务总线实现
│   ├── sconfig/             # 配置管理
│   ├── setcd/               # etcd 客户端封装
│   ├── slog/                # 结构化日志系统
│   ├── snet/                # 网络层 (HTTP、WebSocket)
│   ├── srpc/                # RPC 服务框架
│   ├── store/               # 存储抽象层
│   └── utils/               # 工具函数
├── server/                  # 服务接口定义
├── tool/                    # 开发工具
│   └── simplectl/           # CLI 代码生成器
├── proto/                   # Protocol Buffer 定义
├── config/                  # 配置文件
└── deploy/                  # 部署脚本
```

## 技术栈

| 类别 | 技术 |
|----------|------------|
| **开发语言** | Go 1.22+ |
| **Web 框架** | Gin |
| **RPC 框架** | rpcx |
| **数据库** | MySQL (GORM)、MongoDB |
| **缓存** | Redis |
| **消息队列** | NSQ |
| **服务发现** | etcd |
| **链路追踪** | Jaeger |
| **日志** | Zap |
| **命令行** | Cobra |
| **配置管理** | Viper |

## 快速开始

### 环境要求

- Go 1.22 或更高版本
- etcd 3.4+
- MySQL 8.0+ (可选)
- Redis 6.0+ (可选)

### 安装

```bash
# 安装 CLI 工具
go install github.com/wwengg/simple/tool/simplectl@latest

# 初始化新项目
simplectl rpc init --author "你的名字 <email@example.com>"

# 或者使用 go get 添加到你的项目
go get github.com/wwengg/simple
```

### 创建第一个服务

```bash
# 创建新的 proto 文件
simplectl proto new user

# 从 proto 生成代码
protoc --proto_path=proto --go_out=proto --simple_out=model pbuser/pbuser.proto

# 运行你的服务
go run server/main.go -c config/config.yaml
```

### 配置

创建 `config.yaml` 配置文件：

```yaml
slog:
  level: debug
  format: json

gateway:
  address: ":8080"
  prefix: "/api/v1"

rpc:
  address: ":9090"

etcd:
  endpoints:
    - "127.0.0.1:2379"

redis:
  address: "127.0.0.1:6379"
  db: 0

mysql:
  host: "127.0.0.1"
  port: 3306
  database: "mydb"
```

## 核心模块

### 服务总线 (sbus)

服务总线处理服务间通信，支持多种协议：

```go
// 创建新的服务总线
bus := sbus.NewServiceBus(
    sbus.WithProtocol("quic"),
    sbus.WithHeartbeat(30*time.Second),
)

// 启动服务总线
bus.Start()
```

### RPC 服务

构建可扩展的 RPC 服务，自动服务发现：

```go
// 创建新的 RPC 服务器
server := srpc.NewServer(
    srpc.WithAddress(":9090"),
    srpc.WithEtcdEndpoints([]string{"127.0.0.1:2379"}),
)

// 注册你的服务
pb.RegisterUserService(server, &UserService{})
```

### HTTP 网关

通过 HTTP 暴露服务，内置中间件支持：

```go
// 创建网关
gateway := snet.NewGateway(
    snet.WithAddress(":8080"),
    snet.WithPrefix("/api/v1"),
    snex.WithMiddleware(authMiddleware, rateLimitMiddleware),
)
```

### 配置管理

从文件或环境变量加载配置：

```go
// 从文件加载配置
config := sconfig.Load("config.yaml")

// 访问配置
logLevel := config.Slog.Level
dbHost := config.MySQL.Host
```

## 文档

- [架构指南](docs/architecture.md)
- [API 参考](docs/api.md)
- [部署指南](docs/deployment.md)
- [贡献指南](CONTRIBUTING.md)

## 示例

查看 [examples](examples/) 目录获取示例实现：

- [基础 RPC 服务](examples/basic-rpc)
- [HTTP 网关](examples/http-gateway)
- [服务发现](examples/service-discovery)
- [WebSocket 聊天](examples/websocket-chat)

## 测试

```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...

# 运行竞态检测
go test -race ./...
```

## 路线图

- [ ] Python SDK 支持
- [ ] Kubernetes Operator
- [ ] GraphQL 网关
- [ ] gRPC 转码
- [ ] 服务网格仪表板
- [ ] 自动 API 文档生成

## 贡献

我们欢迎各种形式的贡献！详情请参阅 [CONTRIBUTING.md](CONTRIBUTING.md)。

## 开源协议

Apache License 2.0 - 详见 [LICENSE](LICENSE) 文件。

## 致谢

- [rpcx](https://github.com/smallnest/rpcx) - 高性能 RPC 框架
- [gin-gonic](https://github.com/gin-gonic/gin) - HTTP Web 框架
- [etcd](https://github.com/etcd-io/etcd) - 分布式键值存储
- [jaeger](https://github.com/jaegertracing/jaeger) - 分布式追踪平台

## 联系方式

- 作者：wwengg
- 邮箱：info@wwengg.cn
- 问题反馈：[GitHub Issues](https://github.com/wwengg/simple/issues)
