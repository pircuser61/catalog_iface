package api

import (
	"context"

	"github.com/pkg/errors"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (api *Implementation) GoodCreate(ctx context.Context, in *pb.GoodCreateRequest) (*emptypb.Empty, error) {
	if err := validate(api, in.GetName(), in.GetUnitOfMeasure(), in.GetCountry()); err != nil {
		return nil, err
	}
	return api.Client.GoodCreate(ctx, in)
}

func (api *Implementation) GoodUpdate(ctx context.Context, in *pb.GoodUpdateRequest) (*emptypb.Empty, error) {
	if err := validate(api, in.Good.GetName(), in.Good.GetUnitOfMeasure(), in.Good.GetCountry()); err != nil {
		return nil, err
	}
	return api.Client.GoodUpdate(ctx, in)
}

func (api *Implementation) GoodDelete(ctx context.Context, in *pb.GoodDeleteRequest) (*emptypb.Empty, error) {
	return api.Client.GoodDelete(ctx, in)
}

func (api *Implementation) GoodList(ctx context.Context, in *pb.GoodListRequest) (*pb.GoodListResponse, error) {
	return api.Client.GoodList(ctx, in)

}

func (api *Implementation) GoodGet(ctx context.Context, in *pb.GoodGetRequest) (*pb.GoodGetResponse, error) {
	return api.Client.GoodGet(ctx, in)
}

func validate(api *Implementation, name, uom, country string) error {
	if len(name) < 3 || len(name) > 40 {
		return errors.Errorf("bad name <%v>", name)
	}

	if len(uom) > 10 {
		return errors.Errorf("bad unit of measure <%v>", uom)
	}

	if len(country) < 3 || len(country) > 20 {
		return errors.Errorf("bad country <%v>", country)
	}
	return nil
}
