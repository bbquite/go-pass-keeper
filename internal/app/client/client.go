package client

import (
	"fmt"
	"google.golang.org/grpc/credentials"

	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	"google.golang.org/grpc"
)

type GRPCClient struct {
	conn      *grpc.ClientConn
	PBService pb.PassKeeperServiceClient
}

func NewGRPCClient(serverAddress string, rootCertPath string) (*GRPCClient, error) {
	TLScreeds, err := credentials.NewClientTLSFromFile(rootCertPath, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load CA certificate: %v", err)
	}

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(TLScreeds))
	if err != nil {
		return nil, fmt.Errorf("error init gRPC client: %v", err)
	}

	pbService := pb.NewPassKeeperServiceClient(conn)

	return &GRPCClient{
		conn:      conn,
		PBService: pbService,
	}, nil
}

func (c *GRPCClient) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close client connection: %v", err)
	}
	return nil
}
