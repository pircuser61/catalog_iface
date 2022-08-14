package api

import (
	"context"
	"sync"

	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Implementation struct {
	pb.UnimplementedCatalogIfaceServer
	conn          *grpc.ClientConn
	catalogClient pb.CatalogClient

	goodGet   pb.Catalog_GoodGetClient
	goodGetMu sync.Mutex

	goodList   pb.Catalog_GoodListClient
	goodListMu sync.Mutex

	goodCreate   pb.Catalog_GoodCreateClient
	goodCreateMu sync.Mutex

	goodUpdate   pb.Catalog_GoodUpdateClient
	goodUpdateMu sync.Mutex

	goodDelete   pb.Catalog_GoodDeleteClient
	goodDeleteMu sync.Mutex
}

func New(ctx context.Context) (pb.CatalogIfaceServer, error) {
	conn, err := grpc.Dial(config.GrpcStoreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewCatalogClient(conn)

	api := Implementation{conn: conn, catalogClient: client}

	api.goodGet, err = client.GoodGet(ctx)
	if err != nil {
		return nil, err
	}

	api.goodList, err = client.GoodList(ctx)
	if err != nil {
		return nil, err
	}

	api.goodCreate, err = client.GoodCreate(ctx)
	if err != nil {
		return nil, err
	}

	api.goodUpdate, err = client.GoodUpdate(ctx)
	if err != nil {
		return nil, err
	}

	api.goodDelete, err = client.GoodDelete(ctx)
	if err != nil {
		return nil, err
	}

	return &api, nil
}
