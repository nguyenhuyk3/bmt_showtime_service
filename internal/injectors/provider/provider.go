package provider

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
	"fmt"
	"log"
	"product"
	"sync"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func ProvidePgxPool() *pgxpool.Pool {
	return global.Postgresql
}

func ProvideQueries() *sqlc.Queries {
	return sqlc.New(global.Postgresql)
}

var (
	productClient product.ProductClient

	GRPCClientOnce sync.Once
)

func ProvideProductClient() product.ProductClient {
	GRPCClientOnce.Do(func() {
		conn, err := grpc.Dial(
			fmt.Sprintf("localhost:%s", global.Config.Server.ProductRPCServerPort),
			grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("cannot connect to product service on port 50033: %v", err)
		}

		productClient = product.NewProductClient(conn)
	})

	return productClient
}
