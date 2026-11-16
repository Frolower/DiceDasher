package roomclient

import (
	"context"
	roomv1 "diceDasher/pkg/room/v1"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	connection *grpc.ClientConn
	rooms      roomv1.RoomServiceClient
}

func New(address string) (*Client, error) {
	connection, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{connection: connection, rooms: roomv1.NewRoomServiceClient(connection)}, nil
}

func (c *Client) Create(ctx context.Context, system string, userID uuid.UUID) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	response, err := c.rooms.CreateRoom(ctx, &roomv1.CreateRoomRequest{System: system, UserId: userID.String()})
	if err != nil {
		return "", err
	}
	return response.GetRoomId(), nil
}

func (c *Client) Close() error { return c.connection.Close() }
