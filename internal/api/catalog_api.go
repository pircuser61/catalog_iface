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

	api.goodCreateMu.Lock()
	defer api.goodCreateMu.Unlock()

	if err := api.goodCreate.Send(in); err != nil {
		return nil, err
	}

	result, err := api.goodCreate.Recv()
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, errors.New(*result.Error)
	}
	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodUpdate(ctx context.Context, in *pb.GoodUpdateRequest) (*emptypb.Empty, error) {
	api.goodUpdateMu.Lock()
	defer api.goodUpdateMu.Unlock()

	if err := api.goodUpdate.Send(in); err != nil {
		return nil, err
	}
	result, err := api.goodUpdate.Recv()
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, errors.New(*result.Error)
	}
	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodDelete(ctx context.Context, in *pb.GoodDeleteRequest) (*emptypb.Empty, error) {
	api.goodDeleteMu.Lock()
	defer api.goodDeleteMu.Unlock()

	if err := api.goodDelete.Send(in); err != nil {
		return nil, err
	}

	result, err := api.goodDelete.Recv()
	if err != nil {
		return nil, err
	}
	if result.Error != nil {
		return nil, errors.New(*result.Error)
	}
	return &emptypb.Empty{}, nil
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

func (api *Implementation) GoodGet(ctx context.Context, in *pb.GoodGetRequest) (*pb.GoodGetResponse, error) {
	api.goodGetMu.Lock()
	defer api.goodGetMu.Unlock()

	if err := api.goodGet.Send(in); err != nil {
		return nil, err
	}
	return api.goodGet.Recv()
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
