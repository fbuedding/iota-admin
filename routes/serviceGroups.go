package routes

import (
	"encoding/json"
	"net/http"

	"github.com/fbuedding/iota-admin/internal/globals"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/fbuedding/iota-admin/web/template"
	"github.com/go-chi/chi/v5"
	"github.com/monoculum/formam/v3"
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

		template.Prepare(r, template.FiwareServices(serviceToServiceGroups, "")).Render(r.Context(), w)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Bad Request"))
			return
		}

		var sg i.ServiceGroup
		var decoder = formam.NewDecoder(&formam.DecoderOptions{TagName: "schema"})
		err = decoder.Decode(r.PostForm, &sg)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}

		log.Debug().Any("PostForm", sg).Send()
	})
	r.Get("/servicePaths", func(w http.ResponseWriter, r *http.Request) {
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
		sgs, err := iota.ListServiceGroups(fs)

		if err != nil {
			log.Err(err).Msg("Could not get fiware services")
			http.Error(w, "Could not get fiware services", http.StatusInternalServerError)
			return
		}
		response := ""
		for _, v := range sgs.Services {
			response = response + "<option value=\"" + v.ServicePath + "\">"
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(response))
	})
	return r
}
func AddServiceGroups(repo fr.FiwareRepo) chi.Router {
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
		fss := services.ToFiwareServices()
		encodedBytes, err := json.Marshal(fss)
		if err != nil {
			log.Err(err).Msg("Could not stringify fiware services")
			http.Error(w, "Could not stringify fiware services", http.StatusInternalServerError)
			return
		}
		template.Prepare(r, template.AddServiceGroupForm(string(encodedBytes))).Render(r.Context(), w)
	})
	return r
}
