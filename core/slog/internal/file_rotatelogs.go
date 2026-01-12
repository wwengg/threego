// @Title
// @Description
// @Author  Wangwengang  2023/12/12 09:38
// @Update  Wangwengang  2023/12/12 09:38
package internal

import (
	"fmt"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/wwengg/threego/core/sconfig"
	"go.uber.org/zap/zapcore"
)

var FileRotatelogs = new(fileRotatelogs)

type fileRotatelogs struct{}

// GetWriteSyncer 获取 zapcore.WriteSyncer

func (r *fileRotatelogs) GetWriteSyncer(level zapcore.Level, config *sconfig.Slog) zapcore.WriteSyncer {
	fileWriter := NewCutter(config.Director, level.String(), config.IsAllInOne, WithCutterFormat("2006-01-02"))
	var writeSyncers []zapcore.WriteSyncer
	writeSyncers = append(writeSyncers, zapcore.AddSync(fileWriter))
	if config.LogInConsole {
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))
	}
	if config.LogInSentry {
		err := sentry.Init(sentry.ClientOptions{
			// Either set your DSN here or set the SENTRY_DSN environment variable.
			Dsn: config.SentryDsn,
			// Enable printing of SDK debug messages.
			// Useful when getting started or trying to figure something out.
			Debug: true,
			//Release: ,
		})
		if err != nil {
			fmt.Printf("sentry init fail, err:%v\n", err)
		}
		level2 := transportLevel(config.LogInSentryLevel)
		if level2 <= level {
			writeSyncers = append(writeSyncers, zapcore.AddSync(&Sentry{level: level}))
		}

	}
	return zapcore.NewMultiWriteSyncer(writeSyncers...)

	//return zapcore.AddSync(fileWriter)
}
