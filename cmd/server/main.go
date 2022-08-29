package main

import (
	"context"

	"github.com/opentracing/opentracing-go"
	logger "gitlab.ozon.dev/pircuser61/catalog_iface/internal/logger"
	svPkg "gitlab.ozon.dev/pircuser61/catalog_iface/internal/server"
	"gitlab.ozon.dev/pircuser61/catalog_iface/internal/tracer"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.Debug("Create tracer")
	tracer, closer, err := tracer.CreateTracer("catalog")
	if err != nil {
		logger.Panic("Create tracer error", zap.Error(err))
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	go svPkg.RunREST(ctx)
	svPkg.RunGRPC(ctx)
}
