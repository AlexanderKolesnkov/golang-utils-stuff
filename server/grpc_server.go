package server

import (
	"context"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/selector"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"os"
)

type GrpcServer struct {
	Srv *grpc.Server
}

func NewGrpcServer(
	component string,
	authFunc auth.AuthFunc,
	creds credentials.TransportCredentials,
) *GrpcServer {
	rpcLogger, logTraceID := newRpcLogger(component)

	allButHealthZ := func(ctx context.Context, callMeta interceptors.CallMeta) bool {
		return healthpb.Health_ServiceDesc.ServiceName != callMeta.Service
	}

	s := &GrpcServer{
		Srv: grpc.NewServer(
			grpc.Creds(creds),
			grpc.ChainUnaryInterceptor(
				logging.UnaryServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
				selector.UnaryServerInterceptor(auth.UnaryServerInterceptor(authFunc), selector.MatchFunc(allButHealthZ)),
			),
			grpc.ChainStreamInterceptor(
				logging.StreamServerInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID)),
				selector.StreamServerInterceptor(auth.StreamServerInterceptor(authFunc), selector.MatchFunc(allButHealthZ)),
			),
		),
	}

	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr))

	return s
}

func (s *GrpcServer) ListenAndServe(network, port string) error {
	addr := fmt.Sprintf(":%s", port)

	reflection.Register(s.Srv)

	listen, err := net.Listen(network, addr)
	if err != nil {
		return fmt.Errorf("net Listen: %v", err)
	}

	if err := s.Srv.Serve(listen); err != nil {
		return fmt.Errorf("serve: %v", err)
	}

	return nil
}

func newRpcLogger(component string) (*slog.Logger, func(ctx context.Context) logging.Fields) {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	rpcLogger := logger.With("service", "gRPC/server", "component", component)
	logTraceID := func(ctx context.Context) logging.Fields {
		if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
			return logging.Fields{"traceID", span.TraceID().String()}
		}
		return nil
	}

	return rpcLogger, logTraceID
}

func interceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}
