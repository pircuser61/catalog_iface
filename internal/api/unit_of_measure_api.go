package api

import (
	"context"

	"github.com/pkg/errors"

	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (api *Implementation) UnitOfMeasureCreate(ctx context.Context, in *pb.UnitOfMeasureCreateRequest) (*emptypb.Empty, error) {
	if err := validateUnitOfMeasure(in.GetName()); err != nil {
		return nil, err
	}
	return api.catalogClient.UnitOfMeasureCreate(ctx, in)
}

func (api *Implementation) UnitOfMeasureUpdate(ctx context.Context, in *pb.UnitOfMeasureUpdateRequest) (*emptypb.Empty, error) {
	if err := validateUnitOfMeasure(in.Unit.GetName()); err != nil {
		return nil, err
	}
	return api.catalogClient.UnitOfMeasureUpdate(ctx, in)
}

func (api *Implementation) UnitOfMeasureDelete(ctx context.Context, in *pb.UnitOfMeasureDeleteRequest) (*emptypb.Empty, error) {
	return api.catalogClient.UnitOfMeasureDelete(ctx, in)
}

func (api *Implementation) UnitOfMeasureList(ctx context.Context, in *emptypb.Empty) (*pb.UnitOfMeasureListResponse, error) {
	return api.catalogClient.UnitOfMeasureList(ctx, in)
}

func (api *Implementation) UnitOfMeasureGet(ctx context.Context, in *pb.UnitOfMeasureGetRequest) (*pb.UnitOfMeasureGetResponse, error) {
	return api.catalogClient.UnitOfMeasureGet(ctx, in)
}

func validateUnitOfMeasure(uom string) error {
	if len(uom) < 2 || len(uom) > 10 {
		return errors.Errorf("bad unit of measure <%v>", uom)
	}
	return nil
}
