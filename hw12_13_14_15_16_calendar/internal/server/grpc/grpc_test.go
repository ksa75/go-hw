package grpcserver_test

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
	pb "mycalendar/api/calendarpb"
	"mycalendar/internal/app"
	grpcserver "mycalendar/internal/server/grpc"
	memorystorage "mycalendar/internal/storage/memory"
)

type testLogger struct{}

func (testLogger) Printf(_ string, _ ...any) {}
func (testLogger) Info(_ string)             {}
func (testLogger) Error(_ string)            {}

func TestIntegration_GRPC_AddGetDeleteEvent(t *testing.T) {
	// --- Setup application
	store := memorystorage.New()
	logger := testLogger{}
	appInstance, err := app.New(logger, store)
	require.NoError(t, err)

	// --- Start gRPC server
	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcSrv := grpc.NewServer()
	server := grpcserver.NewServer(appInstance)
	pb.RegisterCalendarServiceServer(grpcSrv, server)

	go grpcSrv.Serve(lis)
	defer grpcSrv.Stop()

	// --- Connect gRPC client
	conn, err := grpc.NewClient(lis.Addr().String(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarServiceClient(conn)
	ctx := context.Background()

	startAt := time.Now().UTC().Truncate(time.Second)

	// --- AddEvent
	_, err = client.AddEvent(ctx, &pb.EventRequest{
		Event: &pb.Event{
			UserId:       "user1",
			Title:        "Integration test event",
			Description:  "desc",
			StartAt:      timestamppb.New(startAt),
			Duration:     "1h",
			NoticeBefore: 10,
		},
	})
	require.NoError(t, err)

	// --- GetEvents
	getResp, err := client.GetEvents(ctx, &pb.Empty{})
	require.NoError(t, err)
	require.Len(t, getResp.Events, 1)
	require.Equal(t, "Integration test event", getResp.Events[0].Title)

	// --- DeleteEvent
	_, err = client.DeleteEvent(ctx, &pb.DeleteRequest{
		UserId: "user1",
		Start:  timestamppb.New(startAt),
	})
	require.NoError(t, err)

	// --- GetEvents again to confirm deletion
	getResp, err = client.GetEvents(ctx, &pb.Empty{})
	require.NoError(t, err)
	require.Len(t, getResp.Events, 0)
}

func TestIntegration_GRPC_UpdateEvent(t *testing.T) {
	// --- Setup
	store := memorystorage.New()
	logger := testLogger{}
	appInstance, err := app.New(logger, store)
	require.NoError(t, err)

	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcSrv := grpc.NewServer()
	pb.RegisterCalendarServiceServer(grpcSrv, grpcserver.NewServer(appInstance))
	go grpcSrv.Serve(lis)
	defer grpcSrv.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarServiceClient(conn)
	ctx := context.Background()

	startAt := time.Now().UTC().Truncate(time.Second)

	// --- Step 1: AddEvent
	_, err = client.AddEvent(ctx, &pb.EventRequest{
		Event: &pb.Event{
			UserId:       "user1",
			Title:        "Initial Title",
			Description:  "Initial description",
			StartAt:      timestamppb.New(startAt),
			Duration:     "1h",
			NoticeBefore: 10,
		},
	})
	require.NoError(t, err)

	// --- Step 2: UpdateEvent
	_, err = client.UpdateEvent(ctx, &pb.EventRequest{
		Event: &pb.Event{
			UserId:       "user1",
			Title:        "Updated Title",
			Description:  "Updated description",
			StartAt:      timestamppb.New(startAt),
			Duration:     "2h",
			NoticeBefore: 20,
		},
	})
	require.NoError(t, err)

	// --- Step 3: Verify update via GetEvents
	getResp, err := client.GetEvents(ctx, &pb.Empty{})
	require.NoError(t, err)
	require.Len(t, getResp.Events, 1)

	updated := getResp.Events[0]
	require.Equal(t, "Updated Title", updated.Title)
	require.Equal(t, "Updated description", updated.Description)
	require.Equal(t, "2h", updated.Duration)
	require.Equal(t, int32(20), updated.NoticeBefore)

	// --- Cleanup
	_, err = client.DeleteEvent(ctx, &pb.DeleteRequest{
		UserId: "user1",
		Start:  timestamppb.New(startAt),
	})
	require.NoError(t, err)
}

func TestIntegration_GRPC_GetEventsByDayWeekMonth(t *testing.T) {
	// --- Setup
	store := memorystorage.New()
	logger := testLogger{}
	appInstance, err := app.New(logger, store)
	require.NoError(t, err)

	lis, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)

	grpcSrv := grpc.NewServer()
	pb.RegisterCalendarServiceServer(grpcSrv, grpcserver.NewServer(appInstance))
	go grpcSrv.Serve(lis)
	defer grpcSrv.Stop()

	conn, err := grpc.NewClient(lis.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewCalendarServiceClient(conn)
	ctx := context.Background()

	startAt := time.Now().UTC().Truncate(time.Second)

	// --- AddEvent
	_, err = client.AddEvent(ctx, &pb.EventRequest{
		Event: &pb.Event{
			UserId:       "user1",
			Title:        "Test event for range queries",
			Description:  "desc",
			StartAt:      timestamppb.New(startAt),
			Duration:     "2h",
			NoticeBefore: 15,
		},
	})
	require.NoError(t, err)

	// --- GetEventsByDay
	dayResp, err := client.GetEventsByDay(ctx, &pb.DateRequest{Date: timestamppb.New(startAt)})
	require.NoError(t, err)
	require.Len(t, dayResp.Events, 1)
	require.Equal(t, "Test event for range queries", dayResp.Events[0].Title)

	// --- GetEventsByWeek
	weekResp, err := client.GetEventsByWeek(ctx, &pb.DateRequest{Date: timestamppb.New(startAt)})
	require.NoError(t, err)
	require.Len(t, weekResp.Events, 1)
	require.Equal(t, "Test event for range queries", weekResp.Events[0].Title)

	// --- GetEventsByMonth
	monthResp, err := client.GetEventsByMonth(ctx, &pb.DateRequest{Date: timestamppb.New(startAt)})
	require.NoError(t, err)
	require.Len(t, monthResp.Events, 1)
	require.Equal(t, "Test event for range queries", monthResp.Events[0].Title)

	// Clean up
	_, err = client.DeleteEvent(ctx, &pb.DeleteRequest{
		UserId: "user1",
		Start:  timestamppb.New(startAt),
	})
	require.NoError(t, err)
}
