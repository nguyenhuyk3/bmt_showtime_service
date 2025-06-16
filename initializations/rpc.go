package initializations

import (
	"bmt_showtime_service/db/sqlc"
	"bmt_showtime_service/global"
	"bmt_showtime_service/internal/rpc"
	"fmt"
	"log"
	"net"

	rpc_showtime "showtime"

	"google.golang.org/grpc"
)

func initRPC() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", global.Config.Server.RPCServerPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	showtimeRPCServer := rpc.NewShowtimeRPCServer(*sqlc.New(global.Postgresql))
	grpcServer := grpc.NewServer()

	rpc_showtime.RegisterShowtimeServer(grpcServer, showtimeRPCServer)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
