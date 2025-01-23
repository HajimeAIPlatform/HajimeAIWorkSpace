/*
 * @Description:
 * @Author: Devin
 * @Date: 2025-01-23 15:56:32
 */
package raydium

import (
	"context"
	"fmt"
	"time"

	pb "hajime/protos/raydium_service_go_grpc"

	"google.golang.org/grpc"
)

// server is used to implement raydium.SwapServiceServer.
type server struct {
	pb.UnimplementedSwapServiceServer
}

func CallSwap(tokenIn string, tokenOut string, privateKey string, amountIn int64, microLamports int64) (string, error) {
	fmt.Println("Calling Swap RPC...")
	conn, err := grpc.Dial(GetGrpcServerHost(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return "", fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewSwapServiceClient(conn)

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
		return "", fmt.Errorf("could not swap: %v", err)
	}

	if res.Status != 200 {
		return "", fmt.Errorf("swap failed. Status: %d, Message: %s", res.Status, res.Message)
	}

	return res.Data.TxId, nil
}
