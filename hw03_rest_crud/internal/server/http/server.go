package internalhttp

import (
	"context"
	"errors"
	"net/http"

	"github.com/julinserg/julinserg/OtusMicroserviceHomeWork/hw03_rest_crud/internal/app"
)

type Storage interface {
	CreateUser(user app.User) error
	DeleteUser(id string) error
	FindUserById(id string) (app.User, error)
	UpdateUser(user app.User) error
}

type Server struct {
	server   *http.Server
	logger   Logger
	endpoint string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func NewServer(logger Logger, storage Storage, endpoint string) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    endpoint,
		Handler: loggingMiddleware(mux, logger),
	}
	uh := userHandler{logger, storage}
	mux.HandleFunc("/api/v1/", hellowHandler)
	mux.HandleFunc("/api/v1/user/", uh.commonHandler)
	return &Server{server, logger, endpoint}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("http server started on " + s.endpoint)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.logger.Info("http server stopped")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
