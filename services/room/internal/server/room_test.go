package server

import (
	"bytes"
	"context"
	roomv1 "diceDasher/pkg/room/v1"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestCreateRoom(t *testing.T) {
	_, conn, _ := startServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client := roomv1.NewRoomServiceClient(conn)
	var output bytes.Buffer
	previous := log.Writer()
	log.SetOutput(&output)
	defer log.SetOutput(previous)
	userID := uuid.NewString()
	response, err := client.CreateRoom(ctx, &roomv1.CreateRoomRequest{System: "vtmv5", UserId: userID})
	if err != nil {
		t.Fatal(err)
	}
	if response.GetStatus() != "ok" {
		t.Fatalf("status = %q", response.GetStatus())
	}
	roomID, err := uuid.Parse(response.GetRoomId())
	if err != nil || roomID.Version() != 4 {
		t.Fatalf("room_id = %q, want a random UUID v4", response.GetRoomId())
	}
	second, err := client.CreateRoom(ctx, &roomv1.CreateRoomRequest{System: "vtmv5", UserId: userID})
	if err != nil {
		t.Fatal(err)
	}
	if second.GetRoomId() == response.GetRoomId() {
		t.Fatal("each call must generate a new room UUID")
	}
	if !strings.Contains(output.String(), `system="vtmv5"`) || !strings.Contains(output.String(), `user_id="`+userID+`"`) {
		t.Fatalf("missing request fields in logs: %s", output.String())
	}
	for _, request := range []*roomv1.CreateRoomRequest{
		{System: "", UserId: userID},
		{System: "vtmv5", UserId: ""},
		{System: "vtmv5", UserId: "invalid"},
		{System: "vtmv5", UserId: uuid.Nil.String()},
	} {
		if _, err := client.CreateRoom(ctx, request); status.Code(err) != codes.InvalidArgument {
			t.Fatalf("invalid request %v: %v", request, err)
		}
	}
}
