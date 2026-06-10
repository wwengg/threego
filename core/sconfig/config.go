// @Title
// @Description
// @Author  Wangwengang  2023/12/10 12:07
// @Update  Wangwengang  2023/12/10 12:07
package sconfig

var S_CONF Config

type Config struct {
	Slog       Slog            `json:"slog" yaml:"slog" mapstructure:"slog"`
	Gateway    Gateway         `mapstructure:"gateway" json:"gateway" yaml:"gateway"`
	RPC        RPC             `mapstructure:"rpc" yaml:"rpc" json:"rpc"`
	RpcService RpcService      `mapstructure:"rpc-service" yaml:"rpc-service" json:"rpcService"`
	Nsq        Nsq             `mapstructure:"nsq" yaml:"nsq" json:"nsq"`
	Redis      Redis           `mapstructure:"redis" yaml:"redis" json:"redis"`
	RootPath   string          `yaml:"root-path" json:"rootPath" mapstructure:"root-path"`
	DBList     []SpecializedDB `mapstructure:"db-list" json:"db-list" yaml:"db-list"`
	Jaeger     Jaeger          `mapstructure:"jaeger" json:"jaeger" yaml:"jaeger"`
	CertPath   string          `mapstructure:"cert-path" json:"certPath" yaml:"cert-path"`
	KeyPath    string          `mapstructure:"key-path" json:"keyPath" yaml:"key-path"`
	Mongo               Mongo     `mapstructure:"mongo" json:"mongo" yaml:"mongo"`
	RecordIgnoreMethods []string  `mapstructure:"record_ignore_methods" json:"record_ignore_methods" yaml:"record_ignore_methods"`
}
