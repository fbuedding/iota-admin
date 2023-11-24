package routes

import (
	"net/http"

	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/web/template"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

type FiwareServicePost struct {
	Name string `schema:"name,required"`
}

func FiwareService(repo fr.FiwareRepo) chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		services, err := repo.ListFiwareServices()

		if err != nil {
			log.Error().Err(err).Msg("Could get fiware services")
			http.Error(w, "Could get fiware services", http.StatusInternalServerError)
			return
		}
		template.Services(services).Render(r.Context(), w)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		var decoder = schema.NewDecoder()
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Bad Request"))
			return
		}

		var fsp FiwareServicePost
		err = decoder.Decode(&fsp, r.PostForm)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Bad Request"))
			return
		}
		if fsp.Name == "" {
			w.WriteHeader(400)
			w.Write([]byte("Bad Request"))
			return
		}

		err = repo.AddFiwareService(fsp.Name)
		if err != nil {
			log.Error().Err(err).Msg("Could add fiware services")
			http.Error(w, "Could add fiware services", http.StatusInternalServerError)
			return
		}
		services, err := repo.ListFiwareServices()

		if err != nil {
			log.Error().Err(err).Msg("Could get fiware services")
			http.Error(w, "Could get fiware services", http.StatusInternalServerError)
			return
		}
		template.Services(services).Render(r.Context(), w)
	})

	return r
}
