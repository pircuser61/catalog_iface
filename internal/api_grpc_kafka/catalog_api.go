package api

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/pkg/errors"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	log "gitlab.ozon.dev/pircuser61/catalog_iface/internal/log"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Good struct {
	Code          uint64
	Name          string
	UnitOfMeasure string
	Country       string
}

func (api *Implementation) GoodCreate(ctx context.Context, in *pb.GoodCreateRequest) (*emptypb.Empty, error) {
	if err := validate(api, in.GetName(), in.GetUnitOfMeasure(), in.GetCountry()); err != nil {
		return nil, err
	}
	good := Good{Name: in.GetName(),
		UnitOfMeasure: in.GetUnitOfMeasure(),
		Country:       in.GetCountry()}

	bjs, err := json.Marshal(good)
	if err != nil {
		return nil, err
	}

	part, offset, err := api.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: config.Topic_create,
		Key:   sarama.StringEncoder(good.Name),
		Value: sarama.ByteEncoder(bjs),
	})

	if err != nil {
		return nil, err
	}
	log.Msgf("part: %d topic %s offset %d\n", part, config.Topic_create, offset)

	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodUpdate(ctx context.Context, in *pb.GoodUpdateRequest) (*emptypb.Empty, error) {
	if err := validate(api, in.Good.GetName(), in.Good.GetUnitOfMeasure(), in.Good.GetCountry()); err != nil {
		return nil, err
	}
	good := Good{Code: in.Good.GetCode(),
		Name:          in.Good.GetName(),
		UnitOfMeasure: in.Good.GetUnitOfMeasure(),
		Country:       in.Good.GetCountry()}

	bjs, err := json.Marshal(good)
	if err != nil {
		return nil, err
	}

	part, offset, err := api.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: config.Topic_update,
		Key:   sarama.StringEncoder(good.Name),
		Value: sarama.ByteEncoder(bjs),
	})

	if err != nil {
		return nil, err
	}
	log.Msgf("part: %d topic %s offset %d\n", part, config.Topic_update, offset)

	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodDelete(ctx context.Context, in *pb.GoodDeleteRequest) (*emptypb.Empty, error) {
	code := in.GetCode()

	bjs, err := json.Marshal(code)
	if err != nil {
		return nil, err
	}

	part, offset, err := api.syncProducer.SendMessage(&sarama.ProducerMessage{
		Topic: config.Topic_delete,
		Key:   sarama.StringEncoder(fmt.Sprint(code)),
		Value: sarama.ByteEncoder(bjs),
	})

	if err != nil {
		return nil, err
	}
	log.Msgf("part: %d topic %s offset %d\n", part, config.Topic_delete, offset)
	return &emptypb.Empty{}, nil
}

func (api *Implementation) GoodGet(ctx context.Context, in *pb.GoodGetRequest) (*pb.GoodGetResponse, error) {
	return api.catalogClient.GoodGet(ctx, in)
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
