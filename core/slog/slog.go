// @Title
// @Description
// @Author  Wangwengang  2023/12/12 09:05
// @Update  Wangwengang  2023/12/12 09:05
package slog

import "context"

type Slog interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	DebugX(ctx context.Context, msg string, fields ...Field)
	InfoX(ctx context.Context, msg string, fields ...Field)
	WarnX(ctx context.Context, msg string, fields ...Field)
	ErrorX(ctx context.Context, msg string, fields ...Field)
	Infof(format string, a ...interface{})
	Debugf(format string, a ...interface{})
	Errorf(format string, a ...interface{})
	Warnf(format string, a ...interface{})
	InfoFX(ctx context.Context, format string, a ...interface{})
	ErrorFX(ctx context.Context, format string, a ...interface{})
	DebugFX(ctx context.Context, format string, a ...interface{})
	Fatal(msg string, fields ...Field)
	Fatalf(format string, a ...interface{})
	Panic(msg string, fields ...Field)
	Panicf(format string, a ...interface{})
}
