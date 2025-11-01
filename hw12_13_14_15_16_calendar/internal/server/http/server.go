package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/app"            //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/config"         //nolint:depguard
	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/storage/common" //nolint:depguard
)

type Server struct {
	Address string
	logger  common.LoggerInterface
	app     app.App
	config  config.ServerConfig
	server  *http.Server
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK}
}

func NewServer(app app.App, config config.ServerConfig, logger common.LoggerInterface) *Server {
	address := net.JoinHostPort(config.Host, config.Port)
	mx := http.NewServeMux()
	httpHandler := NewHandler(app, logger)
	mx.HandleFunc("/hello", httpHandler.helloWorldHandler)

	server := &Server{
		Address: address,
		logger:  logger,
		app:     app,
		config:  config,
		server: &http.Server{
			Addr:           address,
			Handler:        loggingMiddleware(mx, logger),
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		},
	}

	http.HandleFunc("/", server.Handler)

	return server
}

func (s *Server) Handler(w http.ResponseWriter, _ *http.Request) {
	_, err := w.Write([]byte("Work"))
	fmt.Printf("%v", err)
}

func (s *Server) Start(_ context.Context) error {
	s.server.Addr = s.Address
	err := s.server.ListenAndServe()
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	<-ctx.Done()
	s.server.Shutdown(ctx)
	return nil
}
