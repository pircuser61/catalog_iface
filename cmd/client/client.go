package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"

	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

func main() {
	conn, err := grpc.Dial(config.GrpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	client := pb.NewCatalogIfaceClient(conn)
	ctx := context.Background()

	var line string
	in := bufio.NewScanner(os.Stdin)

_cmd:
	for {
		fmt.Print("\nCatalog_iface>")

		if !in.Scan() {
			fmt.Println("Scan error")
			continue
		}
		line = in.Text()
		cmd := strings.Split(line, " ")[0]
		switch cmd {
		case "quit":
			fallthrough
		case "q":
			return
		case "list":
			request := pb.GoodListRequest{}
			params := strings.Split(line, " ")
			len := len(params)
			for i := 1; i < len; i += 2 {
				switch params[i] {
				case "limit":
					if i+1 >= len {
						fmt.Println("<limit> must be a number")
						continue _cmd
					}
					limit, err := strconv.ParseUint(params[i+1], 10, 64)
					if err != nil {
						fmt.Println("<limit> must be a number")
						continue _cmd
					}
					request.Limit = limit
				case "offset":
					if i+1 >= len {
						fmt.Println("<offset> must be a number")
						continue _cmd
					}
					offset, err := strconv.ParseUint(params[i+1], 10, 64)
					if err != nil {
						fmt.Println("<offset> must be a number")
						continue _cmd
					}
					request.Offset = offset
				default:
					fmt.Printf("invalid list param <%s>, avail 'limit <int>', 'offset <int>'", params[i])
					continue _cmd
				}
			}

			response, err := client.GoodList(ctx, &request)
			if err == nil {
				for _, good := range response.GetGoods() {
					fmt.Println(good.Code, good.Name)
				}
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "add":
			params := strings.Split(line, " ")
			if len(params) != 4 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			request := pb.GoodCreateRequest{Name: params[1], UnitOfMeasure: params[2], Country: params[3]}
			response, err := client.GoodCreate(ctx, &request)
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "update":
			params := strings.Split(line, " ")
			if len(params) != 5 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			code, err := strconv.ParseUint(params[1], 10, 64)
			if err != nil {
				fmt.Println("<code> must be a number")
				continue
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
				continue
			}
		case "get":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			code, err := strconv.ParseUint(params[1], 10, 64)
			if err != nil {
				fmt.Println("<code> must be a number")
				continue
			}
			response, err := client.GoodGet(ctx, &pb.GoodGetRequest{Code: code})
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "delete":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			code, err := strconv.ParseUint(params[1], 10, 64)
			if err != nil {
				fmt.Println("<code> must be a number")
				continue
			}
			response, err := client.GoodDelete(ctx, &pb.GoodDeleteRequest{Code: code})
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}

		case "listCountry":
			response, err := client.CountryList(ctx, &emptypb.Empty{})
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "addCountry":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			request := pb.CountryCreateRequest{Name: params[1]}
			response, err := client.CountryCreate(ctx, &request)
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "updateCountry":
			params := strings.Split(line, " ")
			if len(params) != 3 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			u64, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				fmt.Println("<country id> must be a number")
				continue
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
				continue
			}
		case "getCountry":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			u64, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				fmt.Println("<country id> must be a number")
				continue
			}
			country_id := uint32(u64)
			response, err := client.CountryGet(ctx, &pb.CountryGetRequest{CountryId: country_id})
			if err == nil {
				fmt.Printf("response: [%v]", response.Country)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "deleteCountry":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			u64, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				fmt.Println("<country id> must be a number")
				continue
			}
			country_id := uint32(u64)
			response, err := client.CountryDelete(ctx, &pb.CountryDeleteRequest{CountryId: country_id})
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}

		case "listUom":
			response, err := client.UnitOfMeasureList(ctx, &emptypb.Empty{})
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "addUom":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			request := pb.UnitOfMeasureCreateRequest{Name: params[1]}
			response, err := client.UnitOfMeasureCreate(ctx, &request)
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "updateUom":
			params := strings.Split(line, " ")
			if len(params) != 3 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			u64, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				fmt.Println("<country id> must be a number")
				continue
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
				continue
			}
		case "getUom":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			u64, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				fmt.Println("<country id> must be a number")
				continue
			}
			unit_of_measure_id := uint32(u64)
			response, err := client.UnitOfMeasureGet(ctx, &pb.UnitOfMeasureGetRequest{UnitOfMeasureId: unit_of_measure_id})
			if err == nil {
				fmt.Printf("response: [%v]", response.Unit)
			} else {
				fmt.Println(err.Error())
				continue
			}
		case "deleteUom":
			params := strings.Split(line, " ")
			if len(params) != 2 {
				fmt.Printf("invalid args %d items <%v>", len(params), params)
				continue
			}
			u64, err := strconv.ParseUint(params[1], 10, 32)
			if err != nil {
				fmt.Println("<country id> must be a number")
				continue
			}
			unit_of_measure_id := uint32(u64)
			response, err := client.UnitOfMeasureDelete(ctx, &pb.UnitOfMeasureDeleteRequest{UnitOfMeasureId: unit_of_measure_id})
			if err == nil {
				fmt.Printf("response: [%v]", response)
			} else {
				fmt.Println(err.Error())
				continue
			}

		default:
			fmt.Printf("Unknown command <%s>\n", line)
		}
	}
}
