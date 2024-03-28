package routes

import (
	"net/http"

	"github.com/fbuedding/iota-admin/web/templates/layouts"
	"github.com/fbuedding/iota-admin/web/templates/pages"
	"github.com/go-chi/chi/v5"
)

func Login() chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		layouts.Login(pages.LoginForm()).Render(r.Context(), w)
	})

	return r
}
