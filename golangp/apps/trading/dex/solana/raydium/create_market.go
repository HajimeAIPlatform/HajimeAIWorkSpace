package raydium

import (
	"context"
	"fmt"
	"hajime/golangp/common/logging"
	"time"

	pb "hajime/protos/raydium_service_go_grpc"

	"google.golang.org/grpc"
)

func CallCreateMarket(privateKey string, mintAAddress string, mintADecimals int, mintBAddress string, mintBDecimals int) (*pb.CreateMarketData, error) {
	logging.Info("Calling CreateToken RPC...")

	// 使用 WithTransportCredentials 代替 WithInsecure
	conn, err := grpc.Dial(GetGrpcServerHost(), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logging.Danger("did not connect: %v", err)
		return nil, fmt.Errorf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewHajimeGrpcServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 180*time.Second)
	defer cancel()

	req := &pb.CreateMarketRequest{
		MintAAddress:  mintAAddress,
		MintADecimals: uint64(mintADecimals),
		MintBAddress:  mintBAddress,
		MintBDecimals: uint64(mintBDecimals),
		PrivateKey:    privateKey,
	}

	res, err := client.CreateMarket(ctx, req)
	if err != nil {
		logging.Warning("could not create market: %v", err)
		return nil, fmt.Errorf("could not create market: %v", err)
	}

	if res.Status != 200 {
		logging.Warning("create market failed. Status: %d, Message: %s", res.Status, res.Message)
		return nil, fmt.Errorf("create market failed. Status: %d, Message: %s", res.Status, res.Message)
	}

	logging.Info("CreateMarket marketId: %s", res.Data.MarketId)
	logging.Info("CreateMarket txIds: %s", res.Data.TxIds)

	return res.Data, nil
}
