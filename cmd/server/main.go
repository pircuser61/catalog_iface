package main

import (
	"context"

	"github.com/opentracing/opentracing-go"

	svPkg "gitlab.ozon.dev/pircuser61/catalog_iface/internal/server"
	"gitlab.ozon.dev/pircuser61/catalog_iface/internal/tracer"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	tracer, closer, err := tracer.CreateTracer("catalog")
	if err != nil {
		panic(err.Error())
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	go svPkg.RunREST(ctx)
	svPkg.RunGRPC(ctx)
}
