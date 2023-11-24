package server

import (
	"fmt"
	"net/http"

	"github.com/fbuedding/iota-admin/internal/pkg/auth"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/internal/pkg/sessionStore"
	"github.com/fbuedding/iota-admin/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var ()

type Server struct {
	Authenticator auth.Authenticator
	SessionStore  sessionStore.SessionStore
	Port          int
	R             chi.Router
}

func New(a auth.Authenticator, st sessionStore.SessionStore, repo fr.FiwareRepo, port int) *Server {

	var s Server

	s.Authenticator = a
	s.SessionStore = st
	s.Port = port
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	//Public Routes
	r.Group(func(r chi.Router) {
		r.Mount("/login", routes.Login())
		r.Mount("/auth", routes.Auth(s.Authenticator, s.SessionStore))
		r.Mount("/assets", routes.StaticAssets())
	})

	//Private Routes, require authentication
	r.Group(func(r chi.Router) {
		r.Use(routes.AuthMiddleware(s.SessionStore))
		r.Mount("/", routes.Index())
		r.Mount("/fiwareService", routes.FiwareService(repo))
	})

	s.R = r

	return &s
}

func (s Server) Start() error {
	fmt.Printf("Server listening on %d \n", s.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.R)
}
