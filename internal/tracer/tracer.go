package tracer

import (
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	configPkg "gitlab.ozon.dev/pircuser61/catalog_iface/config"
)

func CreateTracer(serviceName string) (opentracing.Tracer, io.Closer, error) {
	rc := &jaegerConfig.ReporterConfig{
		LocalAgentHostPort: configPkg.JaegerHostPort,
		LogSpans:           false,
	}

	sc := &jaegerConfig.SamplerConfig{
		Type:  "const",
		Param: 1,
	}

	cfg := jaegerConfig.Configuration{
		ServiceName: serviceName,
		Disabled:    false,
		Reporter:    rc,
		Sampler:     sc,
	}
	return cfg.NewTracer(jaegerConfig.Logger(jaeger.StdLogger))
}
