package grpc

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"github.com/hw12_13_14_15_calendar/api/eventpb"
	"github.com/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
)

type Conf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Logger interface {
	Info(msg string)
	Error(msg string)
}

type Server struct {
	eventpb.UnimplementedCalendarServiceServer
	log        Logger
	cfg        Conf
	app        server.Application
	grpc       *grpc.Server
	unaryChain grpc.ServerOption
}

func NewServer(log Logger, cfg Conf, app server.Application, unaryChain grpc.ServerOption) *Server {
	return &Server{log: log, cfg: cfg, app: app, unaryChain: unaryChain}
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	s.grpc = grpc.NewServer(s.unaryChain)

	eventpb.RegisterCalendarServiceServer(s.grpc, s)

	go func() {
		s.log.Info(fmt.Sprintf("grpc server started. Addr : %s", addr))
		if err := s.grpc.Serve(lis); err != nil {
			s.log.Error("grpc server serve failed")
		}
	}()
	select {
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop() error {
	if s.grpc != nil {
		s.grpc.GracefulStop()
	}
	return nil
}
