package routes

import (
	"net/http"

	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/web/templates"
	"github.com/fbuedding/iota-admin/web/templates/fiware/iotAgent/agent"
	"github.com/go-chi/chi/v5"
	"github.com/gorilla/schema"
	"github.com/rs/zerolog/log"
)

type ioTAgentPost struct {
	Host  string `schema:"host,required"`
	Port  int    `schema:"port,required"`
	Alias string `schema:"alias"`
}

func IoTAgents(repo fr.FiwareRepo) chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var iotaRows fr.IotaRows
		var err error
		name := r.URL.Query().Get("name")
		if name == "" {
			iotaRows, err = repo.ListIotas()
		} else {
			// TODO implement list by name?
			iotaRows, err = repo.ListIotas()
		}

		if err != nil {
			log.Error().Err(err).Msg("Could not list IoT-Agents")
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			return
		}
		agent.IoTAs(iotaRows).Render(r.Context(), w)
	})

	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		decoder := schema.NewDecoder()
		err := r.ParseForm()
		if err != nil {
			log.Error().Err(err).Msg("Could not parse form")
			templates.HandleError(r.Context(), w, err, http.StatusBadRequest)
			return
		}

		var iotaPost ioTAgentPost
		err = decoder.Decode(&iotaPost, r.PostForm)
		if err != nil {
			log.Error().Err(err).Msg("Could not decode into struct")
			templates.HandleError(r.Context(), w, err, http.StatusBadRequest)
			return
		}

		err = repo.AddIota(iotaPost.Host, iotaPost.Port, iotaPost.Alias)
		if err != nil {
			log.Error().Err(err).Msg("Could add fiware services")
			http.Error(w, "Could add fiware services", http.StatusInternalServerError)
			return
		}
		iotaRows, err := repo.ListIotas()
		if err != nil {
			log.Error().Err(err).Msg("Could not list IoT-Agents")
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			return
		}
		agent.IoTAs(iotaRows).Render(r.Context(), w)
	})

	r.Delete("/{id}", func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		err := repo.DeleteIota(id)
		if err != nil {
			switch err {
			case fr.ErrNotFound:
				http.Error(w, "Service not fount", http.StatusNotFound)
			default:
				log.Error().Err(err).Msgf("Could not delte IoT-Agent with id: %s", id)
				templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
	})

	return r
}
