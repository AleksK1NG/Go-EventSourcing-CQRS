package tracing

import (
	"context"
	"encoding/json"
	"github.com/AleksK1NG/es-microservice/pkg/es"
	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"google.golang.org/grpc/metadata"
)

func StartHttpServerTracerSpan(c echo.Context, operationName string) (context.Context, opentracing.Span) {
	spanCtx, err := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(c.Request().Header))
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctx := opentracing.ContextWithSpan(c.Request().Context(), serverSpan)
		return ctx, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(spanCtx))
	ctx := opentracing.ContextWithSpan(c.Request().Context(), serverSpan)

	return ctx, serverSpan
}

func GetTextMapCarrierFromEvent(event es.Event) opentracing.TextMapCarrier {
	metadataMap := make(opentracing.TextMapCarrier)
	err := json.Unmarshal(event.GetMetadata(), &metadataMap)
	if err != nil {
		return metadataMap
	}
	return metadataMap
}

func StartProjectionTracerSpan(ctx context.Context, operationName string, event es.Event) (context.Context, opentracing.Span) {
	textMapCarrierFromMetaData := GetTextMapCarrierFromEvent(event)

	span, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, textMapCarrierFromMetaData)
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		return ctx, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(span))
	ctx = opentracing.ContextWithSpan(ctx, serverSpan)

	return ctx, serverSpan
}

func GetTextMapCarrierFromMetaData(ctx context.Context) opentracing.TextMapCarrier {
	metadataMap := make(opentracing.TextMapCarrier)
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		for key := range md.Copy() {
			metadataMap.Set(key, md.Get(key)[0])
		}
	}
	return metadataMap
}

func StartGrpcServerTracerSpan(ctx context.Context, operationName string) (context.Context, opentracing.Span) {
	textMapCarrierFromMetaData := GetTextMapCarrierFromMetaData(ctx)

	span, err := opentracing.GlobalTracer().Extract(opentracing.TextMap, textMapCarrierFromMetaData)
	if err != nil {
		serverSpan := opentracing.GlobalTracer().StartSpan(operationName)
		ctx = opentracing.ContextWithSpan(ctx, serverSpan)
		return ctx, serverSpan
	}

	serverSpan := opentracing.GlobalTracer().StartSpan(operationName, ext.RPCServerOption(span))
	ctx = opentracing.ContextWithSpan(ctx, serverSpan)

	return ctx, serverSpan
}

func InjectTextMapCarrier(spanCtx opentracing.SpanContext) (opentracing.TextMapCarrier, error) {
	m := make(opentracing.TextMapCarrier)
	if err := opentracing.GlobalTracer().Inject(spanCtx, opentracing.TextMap, m); err != nil {
		return nil, err
	}
	return m, nil
}

func ExtractTextMapCarrier(spanCtx opentracing.SpanContext) opentracing.TextMapCarrier {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return make(opentracing.TextMapCarrier)
	}
	return textMapCarrier
}

func ExtractTextMapCarrierBytes(spanCtx opentracing.SpanContext) []byte {
	textMapCarrier, err := InjectTextMapCarrier(spanCtx)
	if err != nil {
		return []byte("")
	}

	dataBytes, err := json.Marshal(&textMapCarrier)
	if err != nil {
		return []byte("")
	}
	return dataBytes
}

func InjectTextMapCarrierToGrpcMetaData(ctx context.Context, spanCtx opentracing.SpanContext) context.Context {
	if textMapCarrier, err := InjectTextMapCarrier(spanCtx); err == nil {
		md := metadata.New(textMapCarrier)
		ctx = metadata.NewOutgoingContext(ctx, md)
	}
	return ctx
}

func TraceErr(span opentracing.Span, err error) {
	span.SetTag("error", true)
	span.LogKV("error_code", err.Error())
}
