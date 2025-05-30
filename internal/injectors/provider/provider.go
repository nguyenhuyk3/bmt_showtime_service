package provider

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
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
	filmClient     product.ProductClient
	filmClientOnce sync.Once
)

func ProvideFilmClient() product.ProductClient {
	filmClientOnce.Do(func() {
		conn, err := grpc.Dial("localhost:50033", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("cannot connect to product service on port 50033: %v", err)
		}
		filmClient = product.NewProductClient(conn)
	})
	return filmClient
}
