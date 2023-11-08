package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Login() chi.Router {
	r := chi.NewRouter()
  r.Get("/", func(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "web/static/login.html")
	})

	return r
}
