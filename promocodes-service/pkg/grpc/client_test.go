package grpc

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

type mockPromocodeServiceServer struct {
	UnimplementedPromocodeServiceServer
}

func (m *mockPromocodeServiceServer) CreatePromocode(ctx context.Context, in *CreatePromocodeRequest) (*PromocodeResponse, error) {
	return &PromocodeResponse{Id: 1, Code: "TEST123"}, nil
}

func (m *mockPromocodeServiceServer) GetPromocode(ctx context.Context, in *GetPromocodeRequest) (*PromocodeResponse, error) {
	return &PromocodeResponse{Id: in.GetId(), Code: "TEST123"}, nil
}

func (m *mockPromocodeServiceServer) UpdatePromocode(ctx context.Context, in *UpdatePromocodeRequest) (*PromocodeResponse, error) {
	return &PromocodeResponse{Id: in.GetId(), Code: in.GetCode()}, nil
}

func (m *mockPromocodeServiceServer) DeletePromocode(ctx context.Context, in *DeletePromocodeRequest) (*DeletePromocodeResponse, error) {
	return &DeletePromocodeResponse{Success: true}, nil
}

func (m *mockPromocodeServiceServer) ListPromocodes(ctx context.Context, in *ListPromocodesRequest) (*ListPromocodesResponse, error) {
	return &ListPromocodesResponse{
		Promocodes: []*PromocodeResponse{
			{Id: 1, Code: "TEST1"},
			{Id: 2, Code: "TEST2"},
		},
		Total:   2,
		Page:    in.GetPage(),
		PerPage: in.GetPerPage(),
	}, nil
}

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func setupTestEnvironment(t *testing.T) (*grpc.ClientConn, func()) {
	lis = bufconn.Listen(bufSize)
	server := grpc.NewServer()
	RegisterPromocodeServiceServer(server, &mockPromocodeServiceServer{})

	go func() {
		if err := server.Serve(lis); err != nil {
			t.Fatalf("Server exited with error: %v", err)
		}
	}()

	conn, err := grpc.DialContext(context.Background(), "", grpc.WithContextDialer(bufDialer), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	cleanup := func() {
		conn.Close()
		server.Stop()
	}

	return conn, cleanup
}

func TestGRPCGetPromocode(t *testing.T) {
	conn, cleanup := setupTestEnvironment(t)
	defer cleanup()

	client := NewPromocodeServiceClient(conn)

	resp, err := client.GetPromocode(context.Background(), &GetPromocodeRequest{
		Id: 1,
	})

	assert.NoError(t, err)
	assert.Equal(t, "TEST123", resp.Code)
	assert.Equal(t, int64(1), resp.GetId())
}

func TestGRPCUpdatePromocode(t *testing.T) {
	conn, cleanup := setupTestEnvironment(t)
	defer cleanup()

	client := NewPromocodeServiceClient(conn)

	resp, err := client.UpdatePromocode(context.Background(), &UpdatePromocodeRequest{
		Id:   1,
		Code: "UPDATEDCODE",
	})

	assert.NoError(t, err)
	assert.Equal(t, "UPDATEDCODE", resp.Code)
	assert.Equal(t, int64(1), resp.GetId())
}

func TestGRPCDeletePromocode(t *testing.T) {
	conn, cleanup := setupTestEnvironment(t)
	defer cleanup()

	client := NewPromocodeServiceClient(conn)

	resp, err := client.DeletePromocode(context.Background(), &DeletePromocodeRequest{
		Id: 1,
	})

	assert.NoError(t, err)
	assert.True(t, resp.Success)
}

func TestGRPCListPromocodes(t *testing.T) {
	conn, cleanup := setupTestEnvironment(t)
	defer cleanup()

	client := NewPromocodeServiceClient(conn)

	resp, err := client.ListPromocodes(context.Background(), &ListPromocodesRequest{
		Page:    1,
		PerPage: 10,
	})

	assert.NoError(t, err)
	assert.Len(t, resp.Promocodes, 2)
	assert.Equal(t, int32(1), resp.Page)
	assert.Equal(t, int32(10), resp.PerPage)
}
