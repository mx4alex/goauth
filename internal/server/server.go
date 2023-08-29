package server

import (
	"net/http"
	"context"
)

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(HostAddr string, handler http.Handler) error {
	s.httpServer = &http.Server{
		Addr:    HostAddr,
		Handler: handler,
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}