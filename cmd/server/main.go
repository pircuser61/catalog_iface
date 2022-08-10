package main

import (
	"context"

	"gitlab.ozon.dev/pircuser61/catalog_iface/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go server.RunREST(ctx)
	server.RunGRPC()
}
