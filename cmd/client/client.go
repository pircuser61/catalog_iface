package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	otgrpc "github.com/opentracing-contrib/go-grpc"
	"github.com/opentracing/opentracing-go"
	pb "gitlab.ozon.dev/pircuser61/catalog_iface/api"
	config "gitlab.ozon.dev/pircuser61/catalog_iface/config"
	"gitlab.ozon.dev/pircuser61/catalog_iface/internal/tracer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	go cacheListener(ctx)
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
			createCatalog(ctx, line, client, false)
		case "asyncAdd":
			createCatalog(ctx, line, client, true)
		case "update":
			updateCatalog(ctx, line, client, false)
		case "asyncUpdate":
			updateCatalog(ctx, line, client, true)
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
