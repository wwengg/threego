package plugin

import (
	"context"
	"fmt"

	"github.com/smallnest/rpcx/share"
	"github.com/wwengg/threego/core/slog"
	"go.uber.org/zap"
)

type RecordPlugin struct {
	logger slog.Slog
}

func NewRecordPlugin(slog slog.Slog) *RecordPlugin {
	return &RecordPlugin{
		logger: slog,
	}
}

func (p *RecordPlugin) PreCall(ctx context.Context, serviceName, methodName string, args interface{}) (interface{}, error) {
	p.logger.Info(fmt.Sprintf("Start[%s][%s]", serviceName, methodName), zap.Any("args", args))
	return args, nil
}

func (p *RecordPlugin) PostCall(ctx context.Context, serviceName, methodName string, args, reply interface{}, err error) (interface{}, error) {
	md := ctx.Value(share.ResMetaDataKey)
	if md == nil {
		p.logger.Info(fmt.Sprintf("Finish[%s][%s]", serviceName, methodName), zap.Any("reply", reply))
		return reply, nil
	}
	if v, ok := md.(map[string]string)["record"]; ok && v == "1" {
		// 允许敏感返回值不记录
		p.logger.Info(fmt.Sprintf("Finish[%s][%s]", serviceName, methodName), zap.Any("reply", reply))
	}
	return reply, nil
}
