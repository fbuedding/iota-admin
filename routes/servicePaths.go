package routes

import (
	"net/http"

	"github.com/fbuedding/iota-admin/internal/helpers"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/web/templates"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type servicePathsGetQueryParams struct {
	Service  string `form:"service"`
	IotAgent string `form:"iotAgent"`
}

func ServicePaths(repo fr.FiwareRepo) chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		var params servicePathsGetQueryParams
		err := helpers.Decode(r.URL.Query(), &params)
		if params.Service == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}
		if params.IotAgent == "" {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		iota, err := repo.GetIota(params.IotAgent)
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			log.Error().
				Err(err).
				Msgf("Could not get IoT-Agent for id: %s", params.IotAgent)
			return
		}
		servicePaths, err := iota.GetAllServicePathsForService(params.Service)
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			log.Err(err).Msg("Could not get fiware services")
			return
		}
		response := "<option disabled selected value=\"\">Select service path</option>"
		for _, v := range servicePaths {
			response = response + "<option value=\"" + v + "\">" + v + "</option>"
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
	return r
}
