package raydium

import (
	"context"
	"fmt"
	"hajime/golangp/common/logging"
	"time"

	pb "hajime/protos/raydium_service_go_grpc"

	"google.golang.org/grpc"
)

func CallRevokeMintAuthority(privateKey string, tokenMint string) (*pb.RevokeMintAuthorityData, error) {
	logging.Info("Calling RevokeMintAuthority RPC...")

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

	req := &pb.RevokeMintAuthorityRequest{
		PrivateKey: privateKey,
		TokenMint:  tokenMint,
	}

	res, err := client.RevokeMintAuthority(ctx, req)
	if err != nil {
		logging.Warning("could not revoke mint authority: %v", err)
		return nil, fmt.Errorf("could not revoke mint authority: %v", err)
	}

	if res.Status != 200 {
		logging.Warning("revoke mint authority failed. Status: %d, Message: %s", res.Status, res.Message)
		return nil, fmt.Errorf("revoke mint authority failed. Status: %d, Message: %s", res.Status, res.Message)
	}

	logging.Info("Revoke mint authority TokenMint: %s", res.Data.TokenMint)
	logging.Info("Revoke mint authority TxId: %s", res.Data.TxId)

	return res.Data, nil
}
