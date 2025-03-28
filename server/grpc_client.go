package server

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"github.com/AlexanderKolesnkov/golang-utils-stuff/middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"os"
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

func NewSimpleGrpcClient(
	serverCertPath, clientCertPath, clientKeyPath string,
	username, password string,
	serverHost, serverPort string,
	component string,
) (*grpc.ClientConn, error) {
	creds, err := loadTLSCredentials(serverCertPath, clientCertPath, clientKeyPath)
	if err != nil {
		log.Fatalf("failed to load TLS credentials: %v", err)
	}

	authCred := middleware.NewAuthCredentials(createBasicAuthHeader(username, password))

	return NewGrpcClient(
		component,
		serverHost, serverPort,
		creds,
		authCred,
	)
}

func createBasicAuthHeader(login, password string) string {
	credentials := fmt.Sprintf("%s:%s", login, password)
	encoded := base64.StdEncoding.EncodeToString([]byte(credentials))
	return "Basic " + encoded
}

func loadTLSCredentials(certPath, clientCertPath, clientKeyPath string) (credentials.TransportCredentials, error) {
	certPool := x509.NewCertPool()
	cert, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA certificate: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(cert); !ok {
		return nil, fmt.Errorf("failed to append CA certificate")
	}

	clientCert, err := tls.LoadX509KeyPair(clientCertPath, clientKeyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load client certificate: %v", err)
	}

	return credentials.NewTLS(&tls.Config{
		RootCAs:            certPool,
		Certificates:       []tls.Certificate{clientCert},
		InsecureSkipVerify: false, // В идеале, здесь не должно быть InsecureSkipVerify
	}), nil
}
