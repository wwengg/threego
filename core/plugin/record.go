package plugin

import (
	"context"
	"fmt"
	"sync"

	"github.com/smallnest/rpcx/share"
	"github.com/wwengg/threego/core/slog"
	"go.uber.org/zap"
)

type RecordPlugin struct {
	logger      slog.Slog
	ignoreNames map[string]bool
	mu          sync.RWMutex
}

// WithIgnore 创建时指定不记录的 service.method 列表
func WithIgnore(names []string) func(*RecordPlugin) {
	return func(p *RecordPlugin) {
		for _, n := range names {
			p.ignoreNames[n] = true
		}
	}
}

func NewRecordPlugin(slog slog.Slog, opts ...func(*RecordPlugin)) *RecordPlugin {
	p := &RecordPlugin{
		logger:      slog,
		ignoreNames: make(map[string]bool),
	}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func (p *RecordPlugin) isIgnored(serviceName, methodName string) bool {
	key := fmt.Sprintf("%s.%s", serviceName, methodName)
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.ignoreNames[key]
}

// UpdateIgnoreMethods 热更新忽略列表，配置变更时调用
func (p *RecordPlugin) UpdateIgnoreMethods(names []string) {
	newMap := make(map[string]bool, len(names))
	for _, n := range names {
		newMap[n] = true
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.ignoreNames = newMap
}

func GetReqMetaData(ctx context.Context) map[string]string {
	md := ctx.Value(share.ReqMetaDataKey)
	if md == nil {
		return make(map[string]string)
	}
	mapValue := md.(map[string]string)
	return mapValue
}

func GetResMetaData(ctx context.Context) map[string]string {
	md := ctx.Value(share.ResMetaDataKey)
	if md == nil {
		return make(map[string]string)
	}
	mapValue := md.(map[string]string)
	return mapValue
}

func (p *RecordPlugin) PreCall(ctx context.Context, serviceName, methodName string, args interface{}) (interface{}, error) {
	if p.isIgnored(serviceName, methodName) {
		return args, nil
	}
	md := GetReqMetaData(ctx)
	p.logger.InfoX(ctx, fmt.Sprintf("Start[%s][%s]", serviceName, methodName), zap.Any("args", args), zap.Any("reqMetadata", md))
	return args, nil
}

func (p *RecordPlugin) PostCall(ctx context.Context, serviceName, methodName string, args, reply interface{}, err error) (interface{}, error) {
	if p.isIgnored(serviceName, methodName) {
		return reply, nil
	}
	md := GetResMetaData(ctx)
	p.logger.InfoX(ctx, fmt.Sprintf("Finish[%s][%s]", serviceName, methodName), zap.Any("reply", reply), zap.Any("resMetadata", md))
	return reply, nil
}
