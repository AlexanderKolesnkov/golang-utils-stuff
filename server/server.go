package server

import (
	"net/http" // Стандартный пакет Go для работы с HTTP.
	_ "net/http/pprof"
	"time" // Пакет для работы с временем.
)

type Server struct {
	HttpServer *http.Server // HTTP-сервер.
}

func New(port string) *Server {
	return &Server{
		HttpServer: &http.Server{
			Addr:           "0.0.0.0:" + port, // Адрес и порт сервера.
			ReadTimeout:    180 * time.Second, // Таймаут на чтение.
			WriteTimeout:   180 * time.Second, // Таймаут на запись.
			MaxHeaderBytes: 1 << 20,           // Максимальный размер заголовка в байтах.
		},
	}
}

func (s *Server) InitHandler(handler http.Handler) {
	s.HttpServer.Handler = handler
}

func (s *Server) Run() error {
	return s.HttpServer.ListenAndServe() // Запуск HTTP-сервера.
}

func (s *Server) RunTLS(certFile, keyFile string) error {
	return s.HttpServer.ListenAndServeTLS(certFile, keyFile)
}
