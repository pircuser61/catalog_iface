package server

import (
	"context"
	_ "embed"
	_ "expvar"
	"net"
	"net/http"

	"github.com/flowchartsman/swaggerui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	apiPkg "gitlab.ozon.dev/pircuser61/catalog_iface/internal/api"
	apiKafkaPkg "gitlab.ozon.dev/pircuser61/catalog_iface/internal/api_grpc_kafka"
	logger "gitlab.ozon.dev/pircuser61/catalog_iface/internal/logger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed swagger/api.swagger.json
var spec []byte

func RunGRPC(ctx context.Context) {
	listener, err := net.Listen("tcp", config.GrpcAddr)
	if err != nil {
		logger.Panic("net.Listener", zap.Error(err))
	}
	var apiImplementation pb.CatalogIfaceServer
	if config.UseKafka {
		logger.Debug("Using Kafka + gRPC")
		apiImplementation, err = apiKafkaPkg.New(ctx)
	} else {
		logger.Debug("Using gRPC")
		apiImplementation, err = apiPkg.New(ctx)
	}
	if err != nil {
		logger.Panic("API implementation", zap.Error(err))
	}
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(otgrpc.OpenTracingServerInterceptor(opentracing.GlobalTracer())))
	pb.RegisterCatalogIfaceServer(grpcServer, apiImplementation)

	if err = grpcServer.Serve(listener); err != nil {
		logger.Panic("GRPC.Serve", zap.Error(err))
	}
}

func RunREST(ctx context.Context) {
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterCatalogIfaceHandlerFromEndpoint(ctx, gwmux, config.GrpcAddr, opts); err != nil {
		logger.Panic("RunREST grpc client", zap.Error(err))
	}
	mux := http.NewServeMux()
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	mux.Handle("/", gwmux)

	mux.Handle("/debug/", http.DefaultServeMux)

	if err := http.ListenAndServe(config.HttpAddr, mux); err != nil {
		logger.Panic("RunREST http serve", zap.Error(err))
	}
}
