package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	"gitlab.ozon.dev/pircuser61/catalog_iface/internal/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	tracer, closer, err := tracer.CreateTracer("catalog")
	if err != nil {
		panic(err.Error())
	}
	defer closer.Close()
	opentracing.SetGlobalTracer(tracer)

	conn, err := grpc.Dial(config.GrpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(tracer)),
	)
	if err != nil {
		panic(err)
	}
	client := pb.NewCatalogIfaceClient(conn)
	ctx := context.Background()

	go runErrListiner(ctx)
	dlg(ctx, client)
}

func dlg(ctx context.Context, client pb.CatalogIfaceClient) {
	var line string
	in := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("\nCatalog_iface>")

		if !in.Scan() {
			fmt.Println("Scan error")
			continue
		}
		line = in.Text()
		cmd := strings.Split(line, " ")[0]
		switch cmd {
		case "й":
			fallthrough
		case "quit":
			fallthrough
		case "q":
			return

		case "list":
			listCatalog(ctx, line, client)
		case "add":
			createCatalog(ctx, line, client)
		case "update":
			updateCatalog(ctx, line, client)
		case "get":
			getCatalog(ctx, line, client)
		case "delete":
			deleteCatalog(ctx, line, client)

		case "listCountry":
			listCountry(ctx, line, client)
		case "addCountry":
			createCountry(ctx, line, client)
		case "updateCountry":
			updateCountry(ctx, line, client)
		case "getCountry":
			getCountry(ctx, line, client)
		case "deleteCountry":
			deleteCountry(ctx, line, client)

		case "listUom":
			listUnit(ctx, line, client)
		case "addUom":
			createUnit(ctx, line, client)
		case "updateUom":
			updateUnit(ctx, line, client)
		case "getUom":
			getUnit(ctx, line, client)
		case "deleteUom":
			deleteUnit(ctx, line, client)
		default:
			fmt.Printf("Unknown command <%s>\n", line)
		}
	}
}

func createCatalog(ctx context.Context, line string, client pb.CatalogIfaceClient) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "client/good_create")
	span.LogKV("Cmd", line)
	defer span.Finish()

	params := strings.Split(line, " ")
	if len(params) != 4 {
		fmt.Printf("invalid args %d items <%v>", len(params), params)
		return
	}
	request := pb.GoodCreateRequest{Name: params[1], UnitOfMeasure: params[2], Country: params[3]}
	response, err := client.GoodCreate(ctx, &request)
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
	}
}

func updateCatalog(ctx context.Context, line string, client pb.CatalogIfaceClient) {
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
	response, err := client.GoodUpdate(ctx, &request)
	if err == nil {
		fmt.Printf("response: [%v]", response)
	} else {
		fmt.Println(err.Error())
		return
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
