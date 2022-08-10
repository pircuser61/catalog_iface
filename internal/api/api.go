package api

import (
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Implementation struct {
	pb.UnimplementedCatalogIfaceServer
	conn   *grpc.ClientConn
	Client pb.CatalogClient
}

func New() (pb.CatalogIfaceServer, error) {
	conn, err := grpc.Dial(config.GrpcStoreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewCatalogClient(conn)
	return &Implementation{conn: conn, Client: client}, nil
}
