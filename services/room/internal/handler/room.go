package handler

import (
	"context"
	"diceDasher/pkg/logger"
	roomv1 "diceDasher/pkg/room/v1"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	roomv1.UnimplementedRoomServiceServer
}

func New() *Handler { return &Handler{} }

func (h *Handler) CreateRoom(ctx context.Context, request *roomv1.CreateRoomRequest) (*roomv1.CreateRoomResponse, error) {
	if strings.TrimSpace(request.GetSystem()) == "" {
		return nil, status.Error(codes.InvalidArgument, "system is required")
	}
	if id, err := uuid.Parse(request.GetUserId()); err != nil || id == uuid.Nil {
		return nil, status.Error(codes.InvalidArgument, "valid user_id is required")
	}
	logger.Logf(ctx, "CreateRoom system=%q user_id=%q", request.GetSystem(), request.GetUserId())
	return &roomv1.CreateRoomResponse{Status: "ok", RoomId: uuid.NewString()}, nil
}
