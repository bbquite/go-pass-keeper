package client

import (
	"fmt"
	pb "github.com/bbquite/go-pass-keeper/internal/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type GRPCClient struct {
	conn   *grpc.ClientConn
	Client pb.PassKeeperServiceClient
}

func NewGRPCClient(serverAddress string, rootCertPath string) (*GRPCClient, error) {
	TLScreds, err := credentials.NewClientTLSFromFile(rootCertPath, "")
	if err != nil {
		return nil, fmt.Errorf("failed to load CA certificate: %v", err)
	}

	conn, err := grpc.NewClient(serverAddress, grpc.WithTransportCredentials(TLScreds))
	if err != nil {
		return nil, fmt.Errorf("error init gRPC client: %v", err)
	}

	client := pb.NewPassKeeperServiceClient(conn)

	return &GRPCClient{
		conn:   conn,
		Client: client,
	}, nil
}

func (c *GRPCClient) Close() error {
	err := c.conn.Close()
	if err != nil {
		return fmt.Errorf("failed to close client connection: %v", err)
	}
	return nil
}
