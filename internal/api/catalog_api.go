package api

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (api *Implementation) GoodCreate(ctx context.Context, in *pb.GoodCreateRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "iface/good_create")
	defer span.Finish()
	span.LogKV("request", fmt.Sprintf("%s %s %s",
		in.GetName(),
		in.GetUnitOfMeasure(),
		in.GetCountry()))

	if err := validate(api, in.GetName(), in.GetUnitOfMeasure(), in.GetCountry()); err != nil {
		return nil, err
	}
	if _, err := api.catalogClient.GoodCreate(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodUpdate(ctx context.Context, in *pb.GoodUpdateRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "iface/good_update")
	defer span.Finish()
	span.LogKV("request", fmt.Sprintf("%d %s %s %s",
		in.Good.GetCode(),
		in.Good.GetName(),
		in.Good.GetUnitOfMeasure(),
		in.Good.GetCountry()))
	if err := validate(api, in.Good.GetName(), in.Good.GetUnitOfMeasure(), in.Good.GetCountry()); err != nil {
		return nil, err
	}
	if _, err := api.catalogClient.GoodUpdate(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodDelete(ctx context.Context, in *pb.GoodDeleteRequest) (*emptypb.Empty, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "iface/good_delete")
	defer span.Finish()
	span.LogKV("request", fmt.Sprintf("%d", in.GetCode()))
	if _, err := api.catalogClient.GoodDelete(ctx, in); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodGet(ctx context.Context, in *pb.GoodGetRequest) (*pb.GoodGetResponse, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "iface/good_get")
	defer span.Finish()
	span.LogKV("request", fmt.Sprintf("%d", in.GetCode()))
	return api.catalogClient.GoodGet(ctx, in)
}

func (api *Implementation) GoodList(ctx context.Context, in *pb.GoodListRequest) (*pb.GoodListResponse, error) {
	span, _ := opentracing.StartSpanFromContext(ctx, "iface/good_list")
	defer span.Finish()
	span.LogKV("request", fmt.Sprintf("limit: %d offset: %d", in.GetLimit(), in.GetOffset()))

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
