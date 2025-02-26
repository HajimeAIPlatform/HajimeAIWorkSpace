package raydium

import (
	"context"
	"fmt"
	"hajime/golangp/common/logging"
	"time"

	pb "hajime/protos/raydium_service_go_grpc"

	"google.golang.org/grpc"
)

type CreatePoolDataRes struct {
	PoolId   string
	MarketId string
	TxId     string
}

func CallCreatePool(privateKey string, mintAAddress string, mintADecimals int, mintAInitialAmount int64, mintBAddress string, mintBDecimals int, mintBInitialAmount int64, marketId string) (*CreatePoolDataRes, error) {
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

	req := &pb.CreatePoolRequest{
		MintAAddress:       mintAAddress,
		MintADecimals:      uint64(mintADecimals),
		MintAInitialAmount: uint64(mintAInitialAmount),
		MintBAddress:       mintBAddress,
		MintBDecimals:      uint64(mintBDecimals),
		MintBInitialAmount: uint64(mintBInitialAmount),
		PrivateKey:         privateKey,
		MarketId:           marketId,
	}

	res, err := client.CreatePool(ctx, req)
	if err != nil {
		logging.Warning("could not create pool: %v", err)
		return nil, fmt.Errorf("could not create pool: %v", err)
	}

	if res.Status != 200 {
		logging.Warning("create pool failed. Status: %d, Message: %s", res.Status, res.Message)
		return nil, fmt.Errorf("create pool failed. Status: %d, Message: %s", res.Status, res.Message)
	}

	logging.Info("CreatePool poolId: %s", res.Data.PoolId)
	logging.Info("CreatePool MarketId: %s", res.Data.MarketId)
	logging.Info("CreatePool txId: %s", res.Data.TxId)

	resData := &CreatePoolDataRes{
		PoolId:   res.Data.PoolId,
		MarketId: res.Data.MarketId,
		TxId:     res.Data.TxId,
	}

	return resData, nil
}
