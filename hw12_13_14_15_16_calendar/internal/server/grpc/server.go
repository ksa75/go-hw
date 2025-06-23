package grpcserver

import (
	"context"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
	pb "mycalendar/api/calendarpb"
	"mycalendar/internal/storage"
)

type Application interface {
	CreateEvent(ctx context.Context, uID, title, desc, dur string, noticeBefore int32, startAt time.Time) error
	UpdateEvent(ctx context.Context, uID, title, desc, dur string, noticeBefore int32, startAt time.Time) error
	DeleteEvent(ctx context.Context, userID string, start time.Time) error
	GetEvents(ctx context.Context) ([]storage.Event, error)
	GetEventsByDay(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsByWeek(ctx context.Context, date time.Time) ([]storage.Event, error)
	GetEventsByMonth(ctx context.Context, date time.Time) ([]storage.Event, error)
}

type Server struct {
	pb.UnimplementedCalendarServiceServer
	app Application
}

func NewServer(app Application) *Server {
	return &Server{app: app}
}

func (s *Server) AddEvent(ctx context.Context, req *pb.EventRequest) (*pb.Empty, error) {
	e := req.Event
	err := s.app.CreateEvent(ctx, e.UserId, e.Title, e.Description, e.Duration, e.NoticeBefore, e.StartAt.AsTime())
	return &pb.Empty{}, err
}

func (s *Server) UpdateEvent(ctx context.Context, req *pb.EventRequest) (*pb.Empty, error) {
	e := req.Event
	err := s.app.UpdateEvent(ctx, e.UserId, e.Title, e.Description, e.Duration, e.NoticeBefore, e.StartAt.AsTime())
	return &pb.Empty{}, err
}

func (s *Server) DeleteEvent(ctx context.Context, req *pb.DeleteRequest) (*pb.Empty, error) {
	return &pb.Empty{}, s.app.DeleteEvent(ctx, req.UserId, req.Start.AsTime())
}

func (s *Server) GetEvents(ctx context.Context, _ *pb.Empty) (*pb.EventsResponse, error) {
	events, err := s.app.GetEvents(ctx)
	if err != nil {
		return nil, err
	}
	return convertEvents(events), nil
}

func (s *Server) GetEventsByDay(ctx context.Context, req *pb.DateRequest) (*pb.EventsResponse, error) {
	events, err := s.app.GetEventsByDay(ctx, req.Date.AsTime())
	if err != nil {
		return nil, err
	}
	return convertEvents(events), nil
}

func (s *Server) GetEventsByWeek(ctx context.Context, req *pb.DateRequest) (*pb.EventsResponse, error) {
	events, err := s.app.GetEventsByWeek(ctx, req.Date.AsTime())
	if err != nil {
		return nil, err
	}
	return convertEvents(events), nil
}

func (s *Server) GetEventsByMonth(ctx context.Context, req *pb.DateRequest) (*pb.EventsResponse, error) {
	events, err := s.app.GetEventsByMonth(ctx, req.Date.AsTime())
	if err != nil {
		return nil, err
	}
	return convertEvents(events), nil
}

func convertEvents(events []storage.Event) *pb.EventsResponse {
	res := make([]*pb.Event, 0, len(events))
	for _, e := range events {
		res = append(res, &pb.Event{
			UserId:       e.UserID,
			Title:        e.Title,
			Description:  e.Description,
			StartAt:      timestamppb.New(e.StartDateTime),
			Duration:     e.Duration,
			NoticeBefore: e.NoticeBefore,
		})
	}
	return &pb.EventsResponse{Events: res}
}
