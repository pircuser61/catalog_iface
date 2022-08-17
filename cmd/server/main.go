package main

import (
	"context"

	svPkg "gitlab.ozon.dev/pircuser61/catalog_iface/internal/server"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go svPkg.RunREST(ctx)
	svPkg.RunGRPC(ctx)
}
