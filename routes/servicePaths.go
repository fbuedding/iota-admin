package routes

import (
	"net/http"
	"slices"

	"github.com/fbuedding/iota-admin/internal/globals"
	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func ServicePaths() chi.Router {
	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		service := r.URL.Query().Get("service")
		if service == "" {
			http.Error(w, "no service provided", http.StatusUnprocessableEntity)
			return
		}
		fs := i.FiwareService{
			Service:     service,
			ServicePath: "/*",
		}
		iota := i.IoTA{Host: globals.Conf.IoTAHost, Port: globals.Conf.IoTAPort}
		sgs, err := iota.ListConfigGroups(fs)

		if err != nil {
			log.Err(err).Msg("Could not get fiware services")
			http.Error(w, "Could not get fiware services", http.StatusInternalServerError)
			return
		}
		response := ""
		var servicePaths []string
		for _, v := range sgs.Services {
			if !slices.Contains(servicePaths, v.ServicePath) {
				response = response + "<option value=\"" + v.ServicePath + "\">"
				servicePaths = append(servicePaths, v.ServicePath)
			}
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
	return r
}
