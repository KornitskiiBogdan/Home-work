package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
)

type HttpConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Logger interface {
	Info(msg string)
}

type Application interface { // TODO
}

type Server struct {
	log Logger
	cfg HttpConf
	app Application

	server *http.Server
}

func NewServer(logger Logger, cfg HttpConf, app Application) *Server {
	return &Server{
		log: logger,
		cfg: cfg,
		app: app,
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)
	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	s.server = &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(s.log, mux),
	}

	errCh := make(chan error, 1)

	go func() {
		err := s.server.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			errCh <- nil
			return
		}
		errCh <- err
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-errCh:
		return err
	}
}

func (s *Server) Stop(ctx context.Context) error {
	if s.server == nil {
		return nil
	}

	return s.server.Shutdown(ctx)
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world"))
}
