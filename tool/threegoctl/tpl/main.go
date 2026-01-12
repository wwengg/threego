package tpl

import "fmt"

func MainTemplate() []byte {
	return []byte(`/*
{{ .Copyright }}
*/
package main

import "{{ .PkgName }}/cmd"

func main() {
	cmd.Execute()
}
`)
}

func RootTemplate() []byte {
	return []byte(`/*
{{ .Copyright }}
*/
package cmd

import (
{{- if .Viper }}
	"fmt"{{ end }}
	"github.com/smallnest/rpcx/server"
	"github.com/wwengg/threego/core/plugin"
	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/srpc"
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudflare/tableflip"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
	"{{ .PkgName }}/global"
{{- if .Viper }}
	"github.com/spf13/viper"{{ end }}
)

{{ if .Viper -}}
var cfgFile string
{{- end }}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "{{ .AppName }}",
	Short: "A brief description of your application",
	Long: ` + "`" + `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.` + "`" + `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		global.CONFIG.Slog.Prefix = fmt.Sprintf("%s %d", global.CONFIG.Slog.Prefix, os.Getpid())
		global.InitSlog()
		global.InitSRPC()
		global.InitDB(global.LOG)
		global.InitRedis()
	
		// 创建初始化数据库表
		//global.DB_.AutoMigrate(
		//	model.ServerInfo{},
		//)
		
		{{ .AppName }}Serve(global.CONFIG.RPC, global.CONFIG.RpcService)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
{{- if .Viper }}
	cobra.OnInitialize(initConfig)
{{ end }}
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
{{ if .Viper }}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .AppName }}.yaml)")
{{ else }}
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.{{ .AppName }}.yaml)")
{{ end }}
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// {{ .AppName }}Serve starts a server only registers one service.
// You can register more services and only start one server.
// It blocks until the application exits.
func {{ .AppName }}Serve(rpc sconfig.RPC, rpcService sconfig.RpcService) {
	upg, err := tableflip.New(tableflip.Options{
		PIDFile: "",
	})
	if err != nil {
		panic(err)
	}
	defer upg.Stop()

	// Do an upgrade on SIGHUP
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP)
		for range sig {
			err := upg.Upgrade()
			if err != nil {
				global.LOG.Errorf("Upgrade failed: %v", err)
			}
		}
	}()

	s := server.NewServer()
	t, io, err := plugin.NewTracer(服务名称, global.CONFIG.Jaeger.Agent)
	if err == nil {
		defer io.Close()
		p := plugin.NewJaegerPlugin(t)
		s.Plugins.Add(p)
	} else {
		global.LOG.Errorf("NewTracer error,%s", err.Error())
	}

	// 操作记录插件
	recordPlugin := plugin.NewRecordPlugin(global.LOG)
	s.Plugins.Add(recordPlugin)

	// 开启rpcx监控
	//s.EnableProfile = true
	// 关闭rpcxgateway
	s.DisableHTTPGateway = true
	s.DisableJSONRPC = true
	// 服务注册中心
	srpc.AddRegistryPlugin(s, rpc, rpcService)

	// 在此处注册服务
	//s.RegisterName("Example", new(service.Example), "")

	// Listen must be called before Ready
	ln, err := upg.Listen("tcp", fmt.Sprintf("%s:%s", global.CONFIG.RpcService.ServiceAddr, global.CONFIG.RpcService.Port))
	if err != nil {
		global.LOG.Errorf("Can't listen: %v", err)
	}
	go func() {
		global.LOG.Error(s.ServeListener("tcp", ln).Error())
	}()
	pid := os.Getpid()
	global.LOG.Info("ready",zap.Int("pid",pid))
	if err := upg.Ready(); err != nil {
		panic(err)
	}
	<-upg.Exit()
	global.LOG.Info("server Close",zap.Int("pid",pid))
}

{{ if .Viper -}}
// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		// home, err := os.UserHomeDir()
		// cobra.CheckErr(err)

		// Search config in home directory with name ".{{ .AppName }}" (without extension).
		viper.AddConfigPath("./")
		viper.SetConfigType("yaml")
		viper.SetConfigName("{{ .AppName }}")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		viper.WatchConfig()

		viper.OnConfigChange(func(e fsnotify.Event) {
			fmt.Println("config file changed:", e.Name)
			if err := viper.Unmarshal(&global.CONFIG); err != nil {
				fmt.Println(err)
			}
		})

		if err := viper.Unmarshal(&global.CONFIG); err == nil {
			global.CONFIG.Show()
		} else {
			os.Exit(1)
		}
	}
}
{{- end }}
`)
}

func ConfigYamlTemplate(appName string) []byte {
	return []byte(fmt.Sprintf(`slog:
  level: info
  format: console
  director: log
  encode-level: LowercaseColorLevelEncoder
  stacktrace-key: stacktrace
  max-age: 30
  show-line: true
  log-in-console: true
  prefix: %s

rpc-service:
  service-addr: 127.0.0.1
  port: 9001

rpc:
  register-type: etcdv3
  register-addr:
    - 127.0.0.1:23791
    - 127.0.0.1:23792
    - 127.0.0.1:23793
  base-path: local

redis:
  addr: 127.0.0.1:6379
  password:
  db: 0

db-list:
  - disabled: true # 是否启用
    type: mysql # 数据库的类型,目前支持mysql、pgsql
    alias-name: upms # 数据库的名称,注意: alias-name 需要在db-list中唯一
    path: 127.0.0.1
    port: 3306
    config: charset=utf8mb4&parseTime=True&loc=Local
    db-name: upms
    username: root
    password: root
    max-idle-conns: 10
    max-open-conns: 100
    log-mode: error
    log-zap: true


	jaeger:
	agent: 127.0.0.1:6831`, appName))
}
