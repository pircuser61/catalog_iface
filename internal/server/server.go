package server

import (
	"context"
	_ "embed"
	"net"
	"net/http"

	"github.com/flowchartsman/swaggerui"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	apiPkg "gitlab.ozon.dev/pircuser61/catalog_iface/internal/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//go:embed swagger/api.swagger.json
var spec []byte

func RunGRPC(ctx context.Context) {
	listener, err := net.Listen("tcp", config.GrpcAddr)
	if err != nil {
		panic(err)
	}

	apiImplementation, err := apiPkg.New(ctx)
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterCatalogIfaceServer(grpcServer, apiImplementation)

	if err = grpcServer.Serve(listener); err != nil {
		panic(err)
	}
}

func RunREST(ctx context.Context) {
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterCatalogIfaceHandlerFromEndpoint(ctx, gwmux, config.GrpcAddr, opts); err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mux.Handle("/swagger/", http.StripPrefix("/swagger", swaggerui.Handler(spec)))
	mux.Handle("/", gwmux)
	if err := http.ListenAndServe(config.HttpAddr, mux); err != nil {
		panic(err)
	}
}
