package internalgrpc

import (
	"context"
	"net"
	"net/http"
	"syscall"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"    //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config" //nolint:depguard
	logging "github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/grpc/interceptors"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/server/grpc/pb"
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common" //nolint:depguard
	"google.golang.org/grpc"
)

type Server struct {
	Address string
	logger  common.LoggerInterface
	app     app.App
	config  config.GrpcConfig
	server  *grpc.Server
}

type StatusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *StatusResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func NewStatusResponseWriter(w http.ResponseWriter) *StatusResponseWriter {
	return &StatusResponseWriter{w, http.StatusOK}
}

type Service struct {
	pb.UnimplementedEventsServiceServer
	logger common.LoggerInterface
	app    app.App
	config config.GrpcConfig
}

func (s *Service) CreateEvent(_ context.Context, event *pb.Event) (*pb.Event, error) {
	dtoEvent, err := MapperEventToDtoEvent(event)
	if err != nil {
		return nil, err
	}
	dtoEventNew, err := s.app.CreateEvent(dtoEvent)
	if err != nil {
		return nil, err
	}
	return MapperDtoEventToEvent(dtoEventNew)
}

func (s *Service) UpdateEvent(_ context.Context, event *pb.Event) (*pb.Event, error) {
	dtoEvent, err := MapperEventToDtoEvent(event)
	if err != nil {
		return nil, err
	}
	err = s.app.UpdateEvent(dtoEvent)
	if err != nil {
		return nil, err
	}
	return event, nil
}

func (s *Service) GetEvent(_ context.Context, id *pb.ID) (*pb.Event, error) {
	dtoEvent, err := s.app.GetEvent(id.GetUUID())
	if err != nil {
		return nil, err
	}
	return MapperDtoEventToEvent(dtoEvent)
}

func (s *Service) DeleteEvent(_ context.Context, id *pb.ID) (*pb.Response, error) {
	err := s.app.DeleteEvent(id.GetUUID())
	if err != nil {
		return nil, err
	}
	return &pb.Response{Status: 204}, nil
}

func (s *Service) ListEvent(_ context.Context, _ *pb.EmptyRequest) (*pb.Events, error) {
	events, err := s.app.ListEvent()
	if err != nil {
		return nil, err
	}
	pbEvents := pb.Events{Events: make([]*pb.Event, 0, len(events))}
	for _, event := range events {
		pbEvent, err := MapperDtoEventToEvent(event)
		if err != nil {
			return nil, err
		}
		pbEvents.Events = append(pbEvents.Events, pbEvent)
	}
	return &pbEvents, nil
}

func NewServer(app app.App, config config.GrpcConfig, logger common.LoggerInterface) *Server {
	address := net.JoinHostPort(config.Host, config.Port)

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		logging.UnaryServerRequestLoggingInterceptor(logger),
	))

	pb.RegisterEventsServiceServer(server, &Service{
		app:    app,
		logger: logger,
		config: config,
	})

	serverGrpc := &Server{
		Address: address,
		logger:  logger,
		app:     app,
		config:  config,
		server:  server,
	}

	return serverGrpc
}

func (s *Server) Start(ctx context.Context) error {
	lc := net.ListenConfig{
		Control: func(_, _ string, _ syscall.RawConn) error {
			// Optional: Custom control over the raw connection
			return nil
		},
	}
	lsn, err := lc.Listen(ctx, "tcp", s.Address)
	if err != nil {
		return err
	}
	if err := s.server.Serve(lsn); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) {
	<-ctx.Done()
	s.server.GracefulStop()
}
