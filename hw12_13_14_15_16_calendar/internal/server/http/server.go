package internalhttp

import (
	"context"
)

type HttpConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Server struct { // TODO
}

type Logger interface { // TODO
}

type Application interface { // TODO
}

func NewServer(logger Logger, cfg HttpConf, app Application) *Server {
	return &Server{}
}

func (s *Server) Start(ctx context.Context) error {
	// TODO
	<-ctx.Done()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	// TODO
	return nil
}

// TODO
