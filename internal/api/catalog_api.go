package api

import (
	"context"

	"github.com/pkg/errors"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (api *Implementation) GoodCreate(ctx context.Context, in *pb.GoodCreateRequest) (*emptypb.Empty, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	if err := validate(api, in.GetName(), in.GetUnitOfMeasure(), in.GetCountry()); err != nil {
		return nil, err
	}
	var md metadata.MD
	if _, err := api.catalogClient.GoodCreate(ctx, in, grpc.Header(&md)); err != nil {
		return nil, err
	}
	grpc.SendHeader(ctx, md)
	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodUpdate(ctx context.Context, in *pb.GoodUpdateRequest) (*emptypb.Empty, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	if err := validate(api, in.Good.GetName(), in.Good.GetUnitOfMeasure(), in.Good.GetCountry()); err != nil {
		return nil, err
	}

	if _, err := api.catalogClient.GoodUpdate(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodDelete(ctx context.Context, in *pb.GoodDeleteRequest) (*emptypb.Empty, error) {
	if _, err := api.catalogClient.GoodDelete(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodGet(ctx context.Context, in *pb.GoodGetRequest) (*pb.GoodGetResponse, error) {
	result, err := api.catalogClient.GoodGet(ctx, in)
	return result, err
}

func (api *Implementation) GoodList(ctx context.Context, in *pb.GoodListRequest) (*pb.GoodListResponse, error) {
	api.goodListMu.Lock()
	defer api.goodListMu.Unlock()
	if err := api.goodList.Send(in); err != nil {
		return nil, err
	}
	result, err := api.goodList.Recv()
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, errors.New(*result.Error)
	}
	return result, nil
}

func validate(api *Implementation, name, uom, country string) error {
	if len(name) < 3 || len(name) > 40 {
		return errors.Errorf("bad name <%v>", name)
	}

	if err := validateUnitOfMeasure(uom); err != nil {
		return err
	}

	if err := validateCountry(country); err != nil {
		return err
	}
	return nil
}
