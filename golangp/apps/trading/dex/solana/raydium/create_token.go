package raydium

import (
	"context"
	"fmt"
	"hajime/golangp/common/logging"
	"time"

	pb "hajime/protos/raydium_service_go_grpc"

	"google.golang.org/grpc"
)

func CallCreateToken(privateKey string, tokenName string, tokenSymbol string, description string, uri string, tokenSupply int64, tokenDecimals int64) (*pb.CreateTokenData, error) {
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

	req := &pb.CreateTokenRequest{
		PrivateKey: privateKey,
		TokenInfo: &pb.TokenInfo{
			TokenName:     tokenName,
			TokenSymbol:   tokenSymbol,
			Description:   description,
			Uri:           uri,
			TokenSupply:   uint64(tokenSupply),   // 类型转换
			TokenDecimals: uint64(tokenDecimals), // 类型转换
		},
	}

	res, err := client.CreateToken(ctx, req)
	if err != nil {
		logging.Warning("could not create token: %v", err)
		return nil, fmt.Errorf("could not create token: %v", err)
	}

	if res.Status != 200 {
		logging.Warning("create token failed. Status: %d, Message: %s", res.Status, res.Message)
		return nil, fmt.Errorf("create token failed. Status: %d, Message: %s", res.Status, res.Message)
	}

	logging.Info("CreateToken TokenName: %s", res.Data.TokenName)
	logging.Info("CreateToken TokenMint: %s", res.Data.TokenMint)
	logging.Info("CreateToken TxId: %s", res.Data.TxId)
	logging.Info("CreateToken Supply: %d", res.Data.Supply)

	return res.Data, nil
}
