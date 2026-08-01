package internalhttp

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/hw12_13_14_15_calendar/internal/server"
)

const (
	ReadHeaderTimeout = 10 * time.Second
)

type HTTPConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Logger interface {
	Info(msg string)
}

type Server struct {
	log Logger
	cfg HTTPConf
	app server.Application

	server *http.Server
}

func NewServer(logger Logger, cfg HTTPConf, app server.Application) *Server {
	return &Server{
		log: logger,
		cfg: cfg,
		app: app,
	}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /events", s.createEvent)
	mux.HandleFunc("GET /events/{id}", s.getEvent)
	mux.HandleFunc("PUT /events/{id}", s.updateEvent)
	mux.HandleFunc("DELETE /events/{id}", s.deleteEvent)
	mux.HandleFunc("GET /events/day", s.listOnDay)
	mux.HandleFunc("GET /events/week", s.listOnWeek)
	mux.HandleFunc("GET /events/month", s.listOnMonth)
	mux.HandleFunc("/hello", helloHandler)
	return loggingMiddleware(s.log, mux)
}

func (s *Server) Start(ctx context.Context) error {
	addr := net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port))
	s.server = &http.Server{
		Addr:              addr,
		Handler:           s.Handler(),
		ReadHeaderTimeout: ReadHeaderTimeout,
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

func helloHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("hello world"))
}
