package raydium

import (
	"context"
	"fmt"
	"hajime/golangp/common/logging"
	"time"

	pb "hajime/protos/raydium_service_go_grpc"

	"google.golang.org/grpc"
)

// server is used to implement raydium.SwapServiceServer.
type server struct {
	pb.UnimplementedSwapServiceServer
}

func CallSwap(tokenIn string, tokenOut string, privateKey string, amountIn int64, microLamports int64) (string, error) {
	logging.Info("Calling Swap RPC...")
	conn, err := grpc.Dial(GetGrpcServerHost(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logging.Danger("did not connect: %v", err)
		return "", fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSwapServiceClient(conn)

	// ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	req := &pb.SwapRequest{
		TokenIn:       tokenIn,
		TokenOut:      tokenOut,
		PrivateKey:    privateKey,
		AmountIn:      amountIn,
		MicroLamports: microLamports,
	}

	res, err := client.Swap(ctx, req)
	if err != nil {
		logging.Warning("could not swap: %v", err)
		return "", fmt.Errorf("could not swap: %v", err)
	}

	if res.Status != 200 {
		logging.Warning("swap failed. Status: %d, Message: %s", res.Status, res.Message)
		return "", fmt.Errorf("swap failed. Status: %d, Message: %s", res.Status, res.Message)
	}
	logging.Info("Swap TxId: %s", res.Data.TxId)
	return res.Data.TxId, nil
}
