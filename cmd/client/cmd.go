package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

func createCatalog(ctx context.Context, line string, client pb.CatalogIfaceClient, async bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client/good_create")
	span.LogKV("Cmd", line)
	defer span.Finish()

	params := strings.Split(line, " ")
	if len(params) != 4 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	request := pb.GoodCreateRequest{Name: params[1], UnitOfMeasure: params[2], Country: params[3]}
	if async {
		ctx = metadata.AppendToOutgoingContext(ctx, "mode", "async")
		var md metadata.MD
		_, err := client.GoodCreate(ctx, &request,
			grpc.Header(&md),
		)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		tokens := md.Get("token")
		if len(tokens) == 0 {
			fmt.Println("Error: no token in response")
			return
		}
		token := tokens[0]
		if len(token) == 0 {
			fmt.Println("Error: empty token")
			return
		}

		fmt.Println("waiting for result...")
		time.Sleep(time.Second * 3)
		fmt.Printf("get result (token %s)\n", token)
		result := redisClient.Get(token)
		if err := result.Err(); err == redis.Nil {
			fmt.Println("Error: no result in cache")
		} else if err != nil {
			fmt.Println("Error:", err)
		} else {
			bytes, _ := result.Bytes()
			fmt.Println("result:", string(bytes))
		}
	} else {
		response, err := client.GoodCreate(ctx, &request)
		if err == nil {
			fmt.Printf("response: [%v]", response)
		} else {
			fmt.Println(err.Error())
			return
		}
	}
}

func updateCatalog(ctx context.Context, line string, client pb.CatalogIfaceClient, async bool) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client/good_update")
	span.LogKV("Cmd", line)
	defer span.Finish()

	params := strings.Split(line, " ")
	if len(params) != 5 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	code, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		fmt.Println("<code> must be a number")
		return
	}
	request := pb.GoodUpdateRequest{
		Good: &pb.Good{
			Code:          code,
			Name:          params[2],
			UnitOfMeasure: params[3],
			Country:       params[4]}}
	if async {
		ctx = metadata.AppendToOutgoingContext(ctx, "mode", "async")
		_, err := client.GoodUpdate(ctx, &request)
		if err == nil {
			fmt.Println("request sent, result will come later")
		} else {
			fmt.Println(err.Error())
		}
	} else {
		response, err := client.GoodUpdate(ctx, &request)
		if err == nil {
			fmt.Printf("response: [%v]", response)
		} else {
			fmt.Println(err.Error())
		}
	}
}

func getCatalog(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client/good_get")
	span.LogKV("Cmd", line)
	defer span.Finish()

	params := strings.Split(line, " ")
	if len(params) != 2 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	code, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		fmt.Println("<code> must be a number")
		return
	}
	response, err := client.GoodGet(ctx, &pb.GoodGetRequest{Code: code})
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func deleteCatalog(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client/good_delete")
	span.LogKV("Cmd", line)
	defer span.Finish()

	params := strings.Split(line, " ")
	if len(params) != 2 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	code, err := strconv.ParseUint(params[1], 10, 64)
	if err != nil {
		fmt.Println("<code> must be a number")
		return
	}
	response, err := client.GoodDelete(ctx, &pb.GoodDeleteRequest{Code: code})
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func listCatalog(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client/good_list")
	span.LogKV("Cmd", line)
	defer span.Finish()

	request := pb.GoodListRequest{}
	params := strings.Split(line, " ")
	len := len(params)
	for i := 1; i < len; i += 2 {
		switch params[i] {
		case "limit":
			if i+1 >= len {
				fmt.Println("<limit> must be a number")
				return
			}
			limit, err := strconv.ParseUint(params[i+1], 10, 64)
			if err != nil {
				fmt.Println("<limit> must be a number")
				return
			}
			request.Limit = limit
		case "offset":
			if i+1 >= len {
				fmt.Println("<offset> must be a number")
				return
			}
			offset, err := strconv.ParseUint(params[i+1], 10, 64)
			if err != nil {
				fmt.Println("<offset> must be a number")
				return
			}
			request.Offset = offset
		default:
			fmt.Printf("invalid list param <%s>, avail 'limit <int>', 'offset <int>'", params[i])
			return
		}
	}
	response, err := client.GoodList(ctx, &request)
	if err == nil {
		for _, good := range response.GetGoods() {
			fmt.Println(good.Code, good.Name)
		}
	} else {
		fmt.Println(err.Error())
		return
	}
}

func createCountry(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	params := strings.Split(line, " ")
	if len(params) != 2 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	request := pb.CountryCreateRequest{Name: params[1]}
	response, err := client.CountryCreate(ctx, &request)
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func updateCountry(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	params := strings.Split(line, " ")
	if len(params) != 3 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	u64, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		fmt.Println("<country id> must be a number")
		return
	}
	country_id := uint32(u64)

	request := pb.CountryUpdateRequest{
		Country: &pb.Country{
			CountryId: country_id,
			Name:      params[2]}}
	response, err := client.CountryUpdate(ctx, &request)
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func getCountry(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	params := strings.Split(line, " ")
	if len(params) != 2 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	u64, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		fmt.Println("<country id> must be a number")
		return
	}
	country_id := uint32(u64)
	response, err := client.CountryGet(ctx, &pb.CountryGetRequest{CountryId: country_id})
	if err == nil {
		fmt.Printf("response: [%v]", response.Country)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func deleteCountry(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	params := strings.Split(line, " ")
	if len(params) != 2 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	u64, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		fmt.Println("<country id> must be a number")
		return
	}
	country_id := uint32(u64)
	response, err := client.CountryDelete(ctx, &pb.CountryDeleteRequest{CountryId: country_id})
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func listCountry(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	response, err := client.CountryList(ctx, &emptypb.Empty{})
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func createUnit(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	params := strings.Split(line, " ")
	if len(params) != 2 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	request := pb.UnitOfMeasureCreateRequest{Name: params[1]}
	response, err := client.UnitOfMeasureCreate(ctx, &request)
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func updateUnit(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	params := strings.Split(line, " ")
	if len(params) != 3 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	u64, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		fmt.Println("<country id> must be a number")
		return
	}
	unit_of_measure_id := uint32(u64)

	request := pb.UnitOfMeasureUpdateRequest{
		Unit: &pb.UnitOfMeasure{
			UnitOfMeasureId: unit_of_measure_id,
			Name:            params[2]}}
	response, err := client.UnitOfMeasureUpdate(ctx, &request)
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func getUnit(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	params := strings.Split(line, " ")
	if len(params) != 2 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	u64, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		fmt.Println("<country id> must be a number")
		return
	}
	unit_of_measure_id := uint32(u64)
	response, err := client.UnitOfMeasureGet(ctx, &pb.UnitOfMeasureGetRequest{UnitOfMeasureId: unit_of_measure_id})
	if err == nil {
		fmt.Printf("response: [%v]", response.Unit)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func deleteUnit(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	params := strings.Split(line, " ")
	if len(params) != 2 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	u64, err := strconv.ParseUint(params[1], 10, 32)
	if err != nil {
		fmt.Println("<country id> must be a number")
		return
	}
	unit_of_measure_id := uint32(u64)
	response, err := client.UnitOfMeasureDelete(ctx, &pb.UnitOfMeasureDeleteRequest{UnitOfMeasureId: unit_of_measure_id})
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func listUnit(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	response, err := client.UnitOfMeasureList(ctx, &emptypb.Empty{})
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}
