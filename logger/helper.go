package logger

import (
	"context"
	"github.com/bytepowered/flux"
	"github.com/bytepowered/flux/ext"
)

const (
	TraceId = "trace-id"
	Extras  = "extras"
)

func Trace(traceId string) flux.Logger {
	return With(context.WithValue(context.Background(), TraceId, traceId))
}

func TraceWith(traceId string, fields map[string]string) flux.Logger {
	p := context.WithValue(context.Background(), TraceId, traceId)
	return With(context.WithValue(p, Extras, fields))
}

func With(values context.Context) flux.Logger {
	return ext.NewLoggerWith(values)
}

func TraceContext(ctx flux.Context) flux.Logger {
	if nil == ctx {
		return Trace("no-trace-id")
	}
	if ctxLogger, ok := ctx.GetContextLogger(); ok {
		return ctxLogger
	} else {
		return Trace(ctx.RequestId())
	}
}
