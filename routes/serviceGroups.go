package routes

import (
	"net/http"

	"github.com/fbuedding/iota-admin/internal/globals"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/fbuedding/iota-admin/web/template"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

func ServiceGroups(repo fr.FiwareRepo) chi.Router {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {

		var services fr.FiwareServiceRows
		var err error
		fiwareService := r.URL.Query().Get("name")
		if fiwareService == "" {
			services, err = repo.ListFiwareServices()
		} else {
			var service *fr.FiwareServiceRow
			service, err = repo.GetFiwareService(fiwareService)
			services = append(services, *service)
		}
		if err != nil {
			log.Err(err).Msg("Could not get fiware services")
			http.Error(w, "Could not get fiware services", http.StatusInternalServerError)
			return
		}
		iota := i.IoTA{Host: globals.Conf.IoTAHost, Port: globals.Conf.IoTAPort}
		fss := services.ToFiwareServices()
		log.Debug().Int("countServices", len(fss)).Send()
		serviceToServiceGroups := map[string][]i.ServiceGroup{}
		for _, v := range fss {
			sgs, err := iota.ListServiceGroups(*v)

			if err != nil {
				log.Err(err).Msg("Could not get fiware services")
				http.Error(w, "Could not get fiware services", http.StatusInternalServerError)
				return
			}
			log.Debug().Int("countServiceGroups", sgs.Count).Send()
			if sgs.Count != 0 {
				log.Debug().Any("serviceGroups", sgs.Services).Send()
				serviceToServiceGroups[v.Service] = sgs.Services
			}
		}
		// TODO
		template.FiwareServices(serviceToServiceGroups, "").Render(r.Context(), w)
	})
	return r
}
