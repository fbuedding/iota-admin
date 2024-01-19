package server

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/fbuedding/iota-admin/internal/pkg/auth"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/internal/pkg/sessionStore"
	"github.com/fbuedding/iota-admin/routes"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
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
	//r.Use(middleware.Logger)
	r.Use(LoggerMiddleware(&log.Logger))
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
		r.Mount("/configGroups", routes.ConfigGroups(repo))
		r.Mount("/addConfigGroup", routes.AddConfigGroups(repo))
	})

	s.R = r

	return &s
}

func (s Server) Start() error {
	log.Info().Int("Port", s.Port).Msg("Server starting")
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.R)
}
func LoggerMiddleware(logger *zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			log := logger.With().Logger()

			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				t2 := time.Now()

				// Recover and record stack traces in case of a panic
				if rec := recover(); rec != nil {
					log.Error().
						Str("type", "error").
						Timestamp().
						Interface("recover_info", rec).
						Bytes("debug_stack", debug.Stack()).
						Msg("log system error")
					http.Error(ww, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				// log end request
				log.Debug().
					Str("type", "access").
					Timestamp().
					Fields(map[string]interface{}{
						"remote_ip":  r.RemoteAddr,
						"url":        r.URL.Path,
						"proto":      r.Proto,
						"method":     r.Method,
						"user_agent": r.Header.Get("User-Agent"),
						"status":     ww.Status(),
						"latency_ms": float64(t2.Sub(t1).Nanoseconds()) / 1000000.0,
						"bytes_in":   r.Header.Get("Content-Length"),
						"bytes_out":  ww.BytesWritten(),
					}).
					Msg("incoming_request")
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
