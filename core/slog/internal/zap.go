// @Title
// @Description
// @Author  Wangwengang  2023/12/12 09:14
// @Update  Wangwengang  2023/12/12 09:14
package internal

import (
	"fmt"
	"strings"
	"time"

	"github.com/wwengg/threego/core/sconfig"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapEncodeLevel 根据 EncodeLevel 返回 zapcore.LevelEncoder

func zapEncodeLevel(config *sconfig.Slog) zapcore.LevelEncoder {
	switch {
	case config.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		return zapcore.LowercaseLevelEncoder
	case config.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		return zapcore.LowercaseColorLevelEncoder
	case config.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		return zapcore.CapitalLevelEncoder
	case config.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		return zapcore.CapitalColorLevelEncoder
	default:
		return zapcore.LowercaseLevelEncoder
	}
}

// TransportLevel 根据字符串转化为 zapcore.Level

func transportLevel(level string) zapcore.Level {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.DebugLevel
	}
}

// GetEncoder 获取 zapcore.Encoder
func getEncoder(config *sconfig.Slog) zapcore.Encoder {
	if config.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(config))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(config))
}

// GetEncoderConfig 获取zapcore.EncoderConfig

func getEncoderConfig(config *sconfig.Slog) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  config.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapEncodeLevel(config),
		EncodeTime:     customTimeEncoder(config.Prefix),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// GetEncoderCore 获取Encoder的 zapcore.Core
func getEncoderCore(l zapcore.Level, level zap.LevelEnablerFunc, config *sconfig.Slog) zapcore.Core {
	writer := FileRotatelogs.GetWriteSyncer(l, config) // 日志分割
	return zapcore.NewCore(getEncoder(config), writer, level)
}

// CustomTimeEncoder 自定义日志输出时间格式
func customTimeEncoder(prefix string) func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
	return func(t time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(fmt.Sprintf("[%s]", prefix) + t.Format("2006/01/02 - 15:04:05.000"))
	}
}

// GetZapCores 根据配置文件的Level获取 []zapcore.Core
func GetZapCores(config *sconfig.Slog) []zapcore.Core {
	cores := make([]zapcore.Core, 0, 7)
	for level := transportLevel(config.Level); level <= zapcore.FatalLevel; level++ {
		cores = append(cores, getEncoderCore(level, getLevelPriority(level), config))
	}
	return cores
}

// GetLevelPriority 根据 zapcore.Level 获取 zap.LevelEnablerFunc
func getLevelPriority(level zapcore.Level) zap.LevelEnablerFunc {
	switch level {
	case zapcore.DebugLevel:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	case zapcore.InfoLevel:
		return func(level zapcore.Level) bool { // 日志级别
			return level == zap.InfoLevel
		}
	case zapcore.WarnLevel:
		return func(level zapcore.Level) bool { // 警告级别
			return level == zap.WarnLevel
		}
	case zapcore.ErrorLevel:
		return func(level zapcore.Level) bool { // 错误级别
			return level == zap.ErrorLevel
		}
	case zapcore.DPanicLevel:
		return func(level zapcore.Level) bool { // dpanic级别
			return level == zap.DPanicLevel
		}
	case zapcore.PanicLevel:
		return func(level zapcore.Level) bool { // panic级别
			return level == zap.PanicLevel
		}
	case zapcore.FatalLevel:
		return func(level zapcore.Level) bool { // 终止级别
			return level == zap.FatalLevel
		}
	default:
		return func(level zapcore.Level) bool { // 调试级别
			return level == zap.DebugLevel
		}
	}
}
