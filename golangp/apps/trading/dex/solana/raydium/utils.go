package raydium

import (
	"fmt"
	"os"
)

func GetGrpcServerHost() string {
	if os.Getenv("GRPC_SERVER") != "" {
		fmt.Println("GRPC_SERVER: ", os.Getenv("GRPC_SERVER"))
		return os.Getenv("GRPC_SERVER")
	}
	return "localhost:50051"
}
