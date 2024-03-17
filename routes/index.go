package routes

import (
	"net/http"

	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/web/templates"
	"github.com/fbuedding/iota-admin/web/templates/pages"
	"github.com/go-chi/chi/v5"
	log "github.com/rs/zerolog/log"
)

func Index(repo fr.FiwareRepo) chi.Router {
	router := chi.NewRouter()
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/index", http.StatusSeeOther)
	})

	router.Get("/index", func(w http.ResponseWriter, r *http.Request) {
		agents, err := repo.ListIotas()
		if err != nil {
			log.Err(err).Msg("Could not get fiware services")
			templates.HandleError(r.Context(), w, err, 500)
		}

		log.Debug().Any("agents", agents).Send()
		templates.Prepare(r, pages.Index(agents)).Render(r.Context(), w)
		// TODO add multiple IoT-Agent support
	})

	return router
}
