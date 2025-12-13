package logging

import (
	"context"
	"strings"
	"time"

	"github.com/sblazheev/home_work/hw12_13_14_15_calendar/internal/common" //nolint:depguard
	"google.golang.org/grpc"                                                //nolint:depguard
	"google.golang.org/grpc/metadata"                                       //nolint:depguard
	"google.golang.org/grpc/peer"                                           //nolint:depguard
)

func UnaryServerRequestLoggingInterceptor(logger common.LoggerInterface) grpc.UnaryServerInterceptor {
	return func(c context.Context, r interface{}, i *grpc.UnaryServerInfo, h grpc.UnaryHandler) (interface{}, error) {
		userIP := "0.0.0.0"
		userAgent := ""

		if c != nil {
			if peer, ok := peer.FromContext(c); ok {
				userIP = strings.Split(peer.Addr.String(), ":")[0]
			}
			if headers, ok := metadata.FromIncomingContext(c); ok {
				userAgent = headers.Get("user-agent")[0]
			}
		}
		s := time.Now()
		response, err := h(c, r)
		status := 200
		if err != nil {
			status = 503
		}
		l := time.Since(s)

		logger.Info("REQUEST API",
			"data",
			common.LogEntry{
				IP:        userIP,
				Date:      s,
				Path:      i.FullMethod,
				Proto:     "gRPC",
				Method:    strings.Split(i.FullMethod, "/")[2],
				UserAgent: userAgent,
				Status:    status,
				Latency:   int(l.Milliseconds()),
			})
		logger.Debug("gRPC request", "fullmethod", i.FullMethod,
			"Status", status, "Latency", int(l.Milliseconds()),
			"Request", r, "Response", response)
		return response, err
	}
}
