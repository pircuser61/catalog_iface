package api

import (
	"context"

	"github.com/pkg/errors"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (api *Implementation) CountryCreate(ctx context.Context, in *pb.CountryCreateRequest) (*emptypb.Empty, error) {
	if err := validateCountry(in.GetName()); err != nil {
		return nil, err
	}
	return api.catalogClient.CountryCreate(ctx, in)
}

func (api *Implementation) CountryUpdate(ctx context.Context, in *pb.CountryUpdateRequest) (*emptypb.Empty, error) {
	if err := validateCountry(in.Country.GetName()); err != nil {
		return nil, err
	}
	return api.catalogClient.CountryUpdate(ctx, in)
}

func (api *Implementation) CountryDelete(ctx context.Context, in *pb.CountryDeleteRequest) (*emptypb.Empty, error) {
	return api.catalogClient.CountryDelete(ctx, in)
}

func (api *Implementation) CountryList(ctx context.Context, in *emptypb.Empty) (*pb.CountryListResponse, error) {
	return api.catalogClient.CountryList(ctx, in)
}

func (api *Implementation) CountryGet(ctx context.Context, in *pb.CountryGetRequest) (*pb.CountryGetResponse, error) {
	return api.catalogClient.CountryGet(ctx, in)
}

func (api *Implementation) CountryGetByName(ctx context.Context, in *pb.CountryByNameRequest) (*pb.CountryGetResponse, error) {
	return api.catalogClient.CountryGetByName(ctx, in)
}

func validateCountry(country string) error {
	if len(country) < 3 || len(country) > 20 {
		return errors.Errorf("bad country <%v>", country)
	}
	return nil
}
