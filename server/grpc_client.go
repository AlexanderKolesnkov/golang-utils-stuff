package server

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"time"
)

func NewGrpcClient(
	component string,
	host, port string,
	transportCreds credentials.TransportCredentials,
	authCreds credentials.PerRPCCredentials,
) (*grpc.ClientConn, error) {

	rpcLogger, logTraceID := newRpcLogger(component)
	conn, err := grpc.NewClient(
		fmt.Sprintf("%s:%s", host, port),
		grpc.WithTransportCredentials(transportCreds),
		grpc.WithChainUnaryInterceptor(
			timeout.UnaryClientInterceptor(5*time.Second),
			logging.UnaryClientInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID))),
		grpc.WithChainStreamInterceptor(
			logging.StreamClientInterceptor(interceptorLogger(rpcLogger), logging.WithFieldsFromContext(logTraceID))),
		grpc.WithPerRPCCredentials(authCreds),
	)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
