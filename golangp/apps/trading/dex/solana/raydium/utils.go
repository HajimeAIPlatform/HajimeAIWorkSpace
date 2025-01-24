package raydium

import "os"

func GetGrpcServerHost() string {
	if os.Getenv("GRPC_SERVER") != "" {
		return os.Getenv("GRPC_SERVER")
	}
	return "localhost:50051"
}
