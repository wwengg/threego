# Threego

[English](README.md) | [简体中文](README_cn.md)

A high-performance, production-ready microservices framework written in Go. Threego provides a complete toolkit for building scalable, maintainable distributed systems with ease.

## Features

- **Multi-Protocol Support**: HTTP, WebSocket, QUIC, KCP, TCP
- **Service Mesh**: Built-in service discovery, load balancing, and circuit breaking
- **High Performance**: Powered by rpcx framework for efficient RPC communication
- **Developer Friendly**: CLI tool `simplectl` for rapid project scaffolding
- **Observability**: Integrated Jaeger tracing, structured logging, and Sentry error tracking
- **Hot Reload**: Zero-downtime deployment with tableflip
- **Storage Agnostic**: Support for MySQL, Redis, MongoDB out of the box

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                      Gateway Layer                           │
│              (HTTP/WebSocket/QUIC/KCP)                      │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    Service Bus                              │
│            (Message Routing & Dispatch)                     │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                    RPC Services                             │
│              (Business Logic Processing)                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│                   Storage Layer                             │
│           (MySQL / Redis / MongoDB)                         │
└─────────────────────────────────────────────────────────────┘
                       │
┌──────────────────────▼──────────────────────────────────────┐
│              Service Discovery (etcd)                       │
│              Configuration Center                           │
└─────────────────────────────────────────────────────────────┘
```

## Project Structure

```
threego/
├── core/                    # Core framework modules
│   ├── internal/            # Internal constants and types
│   ├── plugin/              # Plugin system (Jaeger, etc.)
│   ├── sbus/                # Service bus implementation
│   ├── sconfig/             # Configuration management
│   ├── setcd/               # etcd client wrapper
│   ├── slog/                # Structured logging system
│   ├── snet/                # Network layer (HTTP, WebSocket)
│   ├── srpc/                # RPC service framework
│   ├── store/               # Storage abstraction layer
│   └── utils/               # Utility functions
├── server/                  # Server interface definitions
├── tool/                    # Development tools
│   └── simplectl/           # CLI code generator
├── proto/                   # Protocol buffer definitions
├── config/                  # Configuration files
└── deploy/                  # Deployment scripts
```

## Technology Stack

| Category | Technology |
|----------|------------|
| **Language** | Go 1.22+ |
| **Web Framework** | Gin |
| **RPC Framework** | rpcx |
| **Database** | MySQL (GORM), MongoDB |
| **Cache** | Redis |
| **Message Queue** | NSQ |
| **Service Discovery** | etcd |
| **Tracing** | Jaeger |
| **Logging** | Zap |
| **CLI** | Cobra |
| **Config** | Viper |

## Quick Start

### Prerequisites

- Go 1.22 or higher
- etcd 3.4+
- MySQL 8.0+ (optional)
- Redis 6.0+ (optional)

### Installation

```bash
# Install the CLI tool
go install github.com/wwengg/simple/tool/simplectl@latest

# Initialize a new project
simplectl rpc init --author "Your Name <email@example.com>"

# Or use go get to add to your project
go get github.com/wwengg/simple
```

### Create Your First Service

```bash
# Create a new proto file
simplectl proto new user

# Generate code from proto
protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative --simple_out=model --simple_opt=paths=source_relative pbuser/pbuser.proto

# Run your service
go run server/main.go -c config/config.yaml
```

### Configuration

Create a `config.yaml` file:

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

## Core Modules

### Service Bus (sbus)

The service bus handles inter-service communication with support for multiple protocols:

```go
// Create a new service bus
bus := sbus.NewServiceBus(
    sbus.WithProtocol("quic"),
    sbus.WithHeartbeat(30*time.Second),
)

// Start the bus
bus.Start()
```

### RPC Service

Build scalable RPC services with automatic service discovery:

```go
// Create a new RPC server
server := srpc.NewServer(
    srpc.WithAddress(":9090"),
    srpc.WithEtcdEndpoints([]string{"127.0.0.1:2379"}),
)

// Register your service
pb.RegisterUserService(server, &UserService{})
```

### HTTP Gateway

Expose your services via HTTP with built-in middleware:

```go
// Create a gateway
gateway := snet.NewGateway(
    snet.WithAddress(":8080"),
    snet.WithPrefix("/api/v1"),
    snex.WithMiddleware(authMiddleware, rateLimitMiddleware),
)
```

### Configuration Management

Load configurations from files or environment variables:

```go
// Load config from file
config := sconfig.Load("config.yaml")

// Access configuration
logLevel := config.Slog.Level
dbHost := config.MySQL.Host
```

## Feature Status

### Gateway
| Feature | Status |
|---------|--------|
| HTTP | ✅ Completed ([IM Example](https://github.com/wwengg/im)) |
| TCP | ✅ Completed ([IM Example](https://github.com/wwengg/im)) |
| KCP | ✅ Completed ([IM Example](https://github.com/wwengg/im)) |
| WebSocket | ✅ Completed ([IM Example](https://github.com/wwengg/im)) |
| QUIC | ✅ Completed |
| Smart Routing | ✅ Completed ([IM Example](https://github.com/wwengg/im)) |
| Rate Limiting | 🚧 In Progress |
| Circuit Breaker | 🚧 Planned |
| Tracing | ✅ Completed |
| Authentication | ✅ Completed |
| Encryption | 🚧 Planned |
| Timeout Control | ✅ Completed |
| Monitoring | 🚧 Planned |
| JSON/Protobuf | 🚧 Planned |

### RPC Service
| Feature | Status |
|---------|--------|
| Project Initialization | ✅ Completed (Powered by [Cobra](https://github.com/spf13/cobra)) |
| Proto Generation | ✅ Completed |
| Model/Service Generation | ✅ Completed (Powered by [GORM](https://github.com/go-gorm/gorm), [rpcx](https://github.com/smallnest/rpcx)) |
| Performance Monitoring | ✅ Completed |
| Logging | ✅ Completed |
| Authentication | 🚧 Planned |
| Tracing | 🚧 Planned |
| Circuit Breaker | 🚧 Planned |
| Java Support | ✅ Completed (Thanks [rpcx-java](https://github.com/smallnest/rpcx-java)) |
| Python Support | 🚧 Planned |

### Middleware
| Feature | Status |
|---------|--------|
| etcd Distributed Lock | ✅ Completed |
| etcd Service Discovery | ✅ Completed |
| Kubernetes Service Discovery | 🚧 Planned |
| NSQ Message Queue | 🚧 Planned |
| Metrics | ✅ Completed |

### Storage
| Feature | Status |
|---------|--------|
| MySQL | ✅ Completed |
| Redis | 🚧 Planned |
| MongoDB | 🚧 Planned |

## Documentation

- [Architecture Guide](docs/architecture.md)
- [API Reference](docs/api.md)
- [Deployment Guide](docs/deployment.md)
- [Contributing Guide](CONTRIBUTING.md)

## Examples

Check out the [examples](examples/) directory for sample implementations:

- [Basic RPC Service](examples/basic-rpc)
- [HTTP Gateway](examples/http-gateway)
- [Service Discovery](examples/service-discovery)
- [WebSocket Chat](examples/websocket-chat)

## Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# Run tests with race detection
go test -race ./...
```

## Roadmap

- [ ] Python SDK support
- [ ] Kubernetes operator
- [ ] GraphQL gateway
- [ ] gRPC transcoding
- [ ] Service mesh dashboard
- [ ] Automatic API documentation generation

## Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## Stargazers over time

[![Stargazers over time](https://starchart.cc/wwengg/simple.svg)](https://starchart.cc/wwengg/simple)

## Acknowledgments

- [rpcx](https://github.com/smallnest/rpcx) - High performance RPC framework
- [gin-gonic](https://github.com/gin-gonic/gin) - HTTP web framework
- [etcd](https://github.com/etcd-io/etcd) - Distributed key-value store
- [jaeger](https://github.com/jaegertracing/jaeger) - Distributed tracing platform
- [cobra](https://github.com/spf13/cobra) - CLI framework

## Contact

- Author: wwengg
- Email: info@wwengg.cn
- Issues: [GitHub Issues](https://github.com/wwengg/simple/issues)
