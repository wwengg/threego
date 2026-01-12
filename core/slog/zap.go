// @Title
// @Description
// @Author  Wangwengang  2023/12/12 09:06
// @Update  Wangwengang  2023/12/12 09:06
package slog

import (
	"context"
	"fmt"
	"os"

	"github.com/wwengg/threego/core/sconfig"
	"github.com/wwengg/threego/core/slog/internal"
	"github.com/wwengg/threego/core/utils"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//type Field = zap.Field

var sLogInstance Slog

type Zap struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	config *sconfig.Slog
}

func NewZapLog(config *sconfig.Slog) *Zap {
	if config == nil {
		panic(fmt.Errorf("请在config.yaml中配置slog \n"))
	}
	if ok, _ := utils.PathExists(config.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", config.Director)
		_ = os.Mkdir(config.Director, os.ModePerm)
	}

	cores := internal.GetZapCores(config)
	logger := zap.New(zapcore.NewTee(cores...))

	if config.ShowLine {
		logger = logger.WithOptions(zap.AddCaller(), zap.AddCallerSkip(1))
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	z := &Zap{
		logger: logger,
		sugar:  sugar,
		config: config,
	}
	setLog(z)
	return z
}

func (z *Zap) Debug(msg string, fields ...Field) {
	z.logger.Debug(msg, fields...)
}

func (z *Zap) Info(msg string, fields ...Field) {
	z.logger.Info(msg, fields...)
}

func (z *Zap) Error(msg string, fields ...Field) {
	z.logger.Error(msg, fields...)
}

func (z *Zap) Warn(msg string, fields ...Field) {
	z.logger.Warn(msg, fields...)
}

func (z *Zap) Infof(format string, a ...interface{}) {
	z.sugar.Infof(format, a...)
}

func (z *Zap) Debugf(format string, a ...interface{}) {
	z.sugar.Debugf(format, a...)
}

func (z *Zap) Errorf(format string, a ...interface{}) {
	z.sugar.Errorf(format, a...)
}

func (z *Zap) Warnf(format string, a ...interface{}) {
	z.sugar.Warnf(format, a...)
}

func (z *Zap) Fatal(msg string, fields ...Field) {
	z.logger.Fatal(msg, fields...)
}
func (z *Zap) Fatalf(format string, a ...interface{}) {
	z.sugar.Fatalf(format, a...)
}

func (z *Zap) Panic(msg string, fields ...Field) {
	z.logger.Panic(msg, fields...)
}
func (z *Zap) Panicf(format string, a ...interface{}) {
	z.sugar.Panicf(format, a...)
}

func (z *Zap) InfoF(format string, v ...interface{}) {
	z.sugar.Infof(format, v...)
}

func (z *Zap) ErrorF(format string, v ...interface{}) {
	z.sugar.Errorf(format, v...)
}

func (z *Zap) DebugF(format string, v ...interface{}) {
	z.sugar.Debugf(format, v...)
}

func (z *Zap) InfoFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	z.sugar.Infof(format, v...)
}

func (z *Zap) ErrorFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	z.sugar.Errorf(format, v...)
}

func (z *Zap) DebugFX(ctx context.Context, format string, v ...interface{}) {
	fmt.Println(ctx)
	z.sugar.Debugf(format, v...)
}

func setLog(slog Slog) {
	sLogInstance = slog
}

func Ins() Slog {
	return sLogInstance
}

func (z *Zap) GetSugaredLogger() *zap.SugaredLogger {
	return z.sugar
}

func (z *Zap) GetLogger() *zap.Logger {
	return z.logger
}
