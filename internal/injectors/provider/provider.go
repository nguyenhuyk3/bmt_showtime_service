package provider

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
	"log"
	"product"

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

func ProvideFilmClient() product.ProductClient {
	conn, err := grpc.Dial("localhost:50033", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("can not connect to 50051: %v", err)
	}

	return product.NewProductClient(conn)
}
