package api

import (
	"context"
	"sync"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Implementation struct {
	pb.UnimplementedCatalogIfaceServer
	conn          *grpc.ClientConn
	catalogClient pb.CatalogClient

	goodList   pb.Catalog_GoodListClient
	goodListMu sync.Mutex
}

func New(ctx context.Context) (pb.CatalogIfaceServer, error) {
	conn, err := grpc.Dial(config.GrpcStoreAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		return nil, err
	}
	client := pb.NewCatalogClient(conn)

	api := Implementation{conn: conn, catalogClient: client}

	api.goodList, err = client.GoodList(ctx)
	if err != nil {
		return nil, err
	}

	return &api, nil
}
