package grpcstore

import (
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type StoreClient struct {
	pb.UnimplementedCatalogIfaceServer
	conn   *grpc.ClientConn
	Client *pb.CatalogClient
}

func New() (*StoreClient, error) {
	conn, err := grpc.Dial(config.GrpcStoreAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	client := pb.NewCatalogClient(conn)
	return &StoreClient{conn: conn, Client: &client}, nil
}

func (client *StoreClient) Close() {
	client.conn.Close()
}
