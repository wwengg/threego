# Copilot instructions for `simple` (简明指导)

## 一句话概览 ✅
- `simple` 是一个面向快速搭建 RPC 服务与 Gateway 的 Go 框架（基于 `rpcx`, `gorm`, `gin`, 支持 `quic/http3`, `nsq`, `etcd`, `jaeger` 等）。

## 先读哪些文件（快速了解代码走向） 🔎
- 配置与启动：`core/sconfig/viper.go`（配置加载、优先级：命令行 `-c` > 环境变量 `SIMPLE_CONFIG` > 默认），`core/internal/constant.go`（常量名，包含 `SIMPLE_CONFIG`、测试/调试配置名）。
- 网关 / HTTP：`server/gateway.go`, `core/snet/http/gin.go`（Gin + QUIC/HTTP3，Public/Private router groups）。
- RPC：`core/srpc/*`（RPC 接口、与 `rpcx` 的集成）。
- 存储：`core/store/mysql.go`（GORM 初始化，注意会 panic 当缺失配置时）。
- 链路追踪：`core/plugin/jaeger.go`；
- 分布式锁/服务发现：`core/setcd/*`（etcd 交互）。
- 工具生成：`tool/simplectl`（项目/ proto / rpc generator，命令实现在 `tool/simplectl/cmd`）。

## 常见开发工作流（具体命令） 🛠️
- 安装 & 代码生成工具：
  - `go install github.com/wwengg/simple/tool/simplectl@latest`（安装 `simplectl`）
  - `go install github.com/wwengg/protoc-gen-simple@latest`（用于 `protoc` 生成 model/service）
- 用 `simplectl` 快速生成：
  - 初始化 RPC 项目：`simplectl rpc init --author "you <you@example.com>"`
  - 生成 proto：`simplectl proto new user`
- `protoc` 示例（项目中 README 有例子）：
  ```bash
  protoc --proto_path=proto --go_out=proto --go_opt=paths=source_relative --simple_out=model --simple_opt=paths=source_relative pbuser/pbuser.proto
  ```
- 本地测试与构建：
  - 单元/集成测试：`go test ./...`（注意：某些包依赖外部服务）
  - 编译：`go build ./...` / `go install ./...`

## 环境与运行要点 ⚠️
- 配置必须被提供：配置文件路径通过 `-c <path>` 或 环境变量 `SIMPLE_CONFIG` 指定（见 `core/sconfig/viper.go`）。如果不提供，`Viper()` 会 panic。
- 常见外部依赖（测试或运行时需要）：
  - MySQL（`core/store`）
  - etcd（`core/setcd`） — 服务发现、分布式锁
  - NSQ（`core/sbus`, `core/snet`）
  - Jaeger/Tracing（`core/plugin`）
  在做集成测试或运行 demo 时，优先使用轻量 mock 或在 CI 中提供测试服务实例。
- Brotli 压缩：已切换为纯 Go 实现 `github.com/andybalholm/brotli`（见 `core/utils`），因此不再需要系统级 brotli 库；Windows 平台仍使用返回原始数据的 stub（见 `core/utils/compressor_win.go`）。
## 项目约定与模式（在本仓库中常见的做法） 📐
- 配置结构集中在 `core/sconfig`，通过 `S_CONF` 全局变量持有解码后的配置信息。
- 日志：自封装的 `Slog` 接口（`core/slog`），底层用 `zap` + 文件切割（`core/slog/internal`），日志格式/级别通过 config 控制。
- 对外网络层采用 `server/gateway.go` + `core/snet/http`，支持同时开启 TCP 与 QUIC（HTTP3）。
- RPC 层与 `rpcx` 紧耦合，存在 `SRPC` 抽象便于替换/测试。

> 重要：某些初始化函数在配置缺失时会直接 `panic`（例如 `GormMysqlByConfig`）；保证测试/开发环境提供最小配置或改用替身/接口注入。

## Examples（代码片段引用） 💡
- 读取配置：见 `core/sconfig/viper.go`（优先级、`-c`、`SIMPLE_CONFIG`）
- GORM 初始化（示例）：`core/store/mysql.go` —— `GormMysqlByConfig` 会用 `m.Dsn()` 构建 DSN 并 `panic` on error
- RPC 拦截器/Tracing：`core/plugin/jaeger.go`（如何从 message metadata 提取/注入 span）

## 什么不写在这里（避免误导） ❗
- 不包含未实现或未在代码中体现的最佳实践（例如：自动化 CI、容器编排方案除非仓库已有配置）。

---

请审阅这份说明：
- 哪些细节需要补充（例如：你们自己的 `config.yaml` 示例、常用 `make` 命令或 CI 步骤）？
- 我可以把你确认的内容合并到仓库并迭代更新。 😊