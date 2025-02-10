package config

import (
	"context"
	"encoding/json"
	"sync"

	"mosn.io/api"
)

const (
	contextKeyStreamID = iota
	contextKeyConnection
	contextKeyConnectionID
	contextKeyConnectionPoolIndex
	contextKeyListenerPort
	contextKeyListenerName
	contextKeyListenerType
	contextKeyListenerStatsNameSpace
	contextKeyNetworkFilterChainFactories
	contextKeyStreamFilterChainFactories
	contextKeyBufferPoolCtx
	contextKeyAccessLogs
	contextOriRemoteAddr
	contextKeyAcceptChan
	contextKeyAcceptBuffer
	contextKeyConnectionFd
	contextKeyTraceSpanKey
	contextKeyActiveSpan
	contextKeyTraceId
	contextKeyVariables
	contextKeyProxyGeneralConfig
	contextKeyDownStreamProtocol
	contextKeyUpStreamProtocol
	contextKeyDownStreamHeaders
	contextKeyDownStreamRespHeaders
	contextKeySharedConnCtx
	contextKeyEnd
)

const ContextKey = "context_key"

func ContextByContext(ctx context.Context) ([]interface{}, bool) {
	cfg, ok := ctx.Value(ContextKey).(*[]interface{})
	return (*cfg), ok
}

func GetSpan(ctx context.Context) (api.Span, bool) {
	cfg, ok := ContextByContext(ctx)
	if !ok {
		return nil, false
	}
	info := cfg[contextKeyTraceSpanKey]
	span, ok := info.(api.Span)
	if !ok {
		return nil, false
	}
	return span, ok
}

//
//func GetProxyGeneralConfig(ctx context.Context) (map[string]interface{}, bool) {
//	cfg, ok := ContextByContext(ctx)
//	if !ok {
//		return nil, false
//	}
//	info := cfg[contextKeyProxyGeneralConfig]
//	if info == nil {
//		return nil, false
//	}
//	if v, ok := info.(map[string]interface{}); ok {
//		return v, ok
//	}
//	return nil, false
//}

func GetSharedConnCtx(ctx context.Context) (*sync.Map, bool) {
	cfg, ok := ContextByContext(ctx)
	if !ok || len(cfg) < contextKeySharedConnCtx {
		return nil, false
	}
	info := cfg[contextKeySharedConnCtx]
	if info == nil {
		return nil, false
	}
	if v, ok := info.(*sync.Map); ok {
		return v, ok
	}
	return nil, false
}

func GetListenerType(ctx context.Context) (string, bool) {
	cfg, ok := ContextByContext(ctx)
	if !ok {
		return "", false
	}
	value := cfg[contextKeyListenerType]
	nv, err := json.Marshal(value)
	if err != nil {
		return "", false
	}
	var data string
	if err := json.Unmarshal(nv, &data); err != nil {
		return "", false
	}
	return data, true
}
