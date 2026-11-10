package handler

import (
	"backend/internal/auth"
	"backend/internal/roomclient"
	"context"
	"diceDasher/pkg/httputil"
	roomv1 "diceDasher/pkg/room/v1"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recordingRoomServer struct {
	roomv1.UnimplementedRoomServiceServer
	requests chan *roomv1.CreateRoomRequest
}

const testRoomID = "9da3cc51-c9a2-4f98-a865-d4aa135c0e96"

func (s *recordingRoomServer) CreateRoom(ctx context.Context, request *roomv1.CreateRoomRequest) (*roomv1.CreateRoomResponse, error) {
	s.requests <- request
	switch request.GetSystem() {
	case "unavailable":
		return nil, status.Error(codes.Unavailable, "private upstream details")
	case "invalid":
		return nil, status.Error(codes.InvalidArgument, "private upstream details")
	case "failure":
		return nil, status.Error(codes.Internal, "private upstream details")
	case "slow":
		<-ctx.Done()
		return nil, status.FromContextError(ctx.Err()).Err()
	}
	return &roomv1.CreateRoomResponse{Status: "ok", RoomId: testRoomID}, nil
}

func TestCreateRoomHTTPToGRPC(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	rpc := grpc.NewServer()
	recorder := &recordingRoomServer{requests: make(chan *roomv1.CreateRoomRequest, 1)}
	roomv1.RegisterRoomServiceServer(rpc, recorder)
	done := make(chan error, 1)
	go func() { done <- rpc.Serve(listener) }()
	t.Cleanup(func() { rpc.Stop(); <-done })
	client, err := roomclient.New(listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	userID := uuid.New()
	authenticator := authenticateFunc{authenticate: func(ctx context.Context, token string) (auth.AccessClaims, error) {
		if token != "valid" {
			return auth.AccessClaims{}, auth.ErrInvalidToken
		}
		return auth.AccessClaims{UserID: userID}, nil
	}}
	router := httputil.NewRouter()
	NewWithOptions(authenticator, nil, Options{Rooms: client}).RegisterRouters(router)
	for _, tc := range []struct {
		name, token, body, system string
		code                      int
	}{
		{"success", "valid", `{"system":"vtmv5"}`, "vtmv5", 200},
		{"another system", "valid", `{"system":"tes"}`, "tes", 200},
		{"no token", "", `{"system":"vtmv5"}`, "", 401},
		{"invalid token", "bad", `{"system":"vtmv5"}`, "", 401},
		{"invalid JSON", "valid", `{`, "", 400},
		{"empty body", "valid", ``, "", 400},
		{"missing system", "valid", `{}`, "", 400},
		{"blank system", "valid", `{"system":"  "}`, "", 400},
		{"null", "valid", `null`, "", 400},
		{"wrong type", "valid", `{"system":123}`, "", 400},
		{"spoof user", "valid", `{"system":"vtmv5","user_id":"someone-else"}`, "", 400},
		{"multiple values", "valid", `{"system":"vtmv5"} {}`, "", 400},
		{"oversized", "valid", `{"system":"` + strings.Repeat("x", 17<<10) + `"}`, "", 413},
		{"unavailable", "valid", `{"system":"unavailable"}`, "unavailable", 503},
		{"invalid RPC", "valid", `{"system":"invalid"}`, "invalid", 400},
		{"RPC failure", "valid", `{"system":"failure"}`, "failure", 502},
		{"RPC timeout", "valid", `{"system":"slow"}`, "slow", 504},
	} {
		t.Run(tc.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodPost, "/room", strings.NewReader(tc.body))
			request.Header.Set("Content-Type", "application/json")
			if tc.token != "" {
				request.Header.Set("Authorization", "Bearer "+tc.token)
			}
			if tc.system == "slow" {
				ctx, cancel := context.WithTimeout(request.Context(), 100*time.Millisecond)
				defer cancel()
				request = request.WithContext(ctx)
			}
			response := httptest.NewRecorder()
			router.ServeHTTP(response, request)
			if response.Code != tc.code {
				t.Fatalf("HTTP = %d, want %d: %s", response.Code, tc.code, response.Body.String())
			}
			if response.Code == http.StatusOK && strings.TrimSpace(response.Body.String()) != `{"status":"ok","room_id":"`+testRoomID+`"}` {
				t.Fatalf("response = %s", response.Body.String())
			}
			if strings.Contains(response.Body.String(), "private upstream details") {
				t.Fatal("upstream details leaked")
			}
			if tc.system != "" {
				select {
				case received := <-recorder.requests:
					if received.GetSystem() != tc.system || received.GetUserId() != userID.String() {
						t.Fatalf("RPC received %v", received)
					}
				default:
					t.Fatal("request did not reach room over gRPC")
				}
			} else {
				select {
				case received := <-recorder.requests:
					t.Fatalf("rejected request reached room: %v", received)
				default:
				}
			}
		})
	}
	rpc.Stop()
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/room", strings.NewReader(`{"system":"vtmv5"}`))
	request.Header.Set("Authorization", "Bearer valid")
	router.ServeHTTP(response, request)
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("stopped room = %d", response.Code)
	}
}
