package routes

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/fbuedding/iota-admin/internal/globals"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/fbuedding/iota-admin/web/templates"
	"github.com/fbuedding/iota-admin/web/templates/components"
	configgroup "github.com/fbuedding/iota-admin/web/templates/fiware/iotAgent/configGroup"
	"github.com/go-chi/chi/v5"
	"github.com/monoculum/formam/v3"
	"github.com/rs/zerolog/log"
)

type CofigGroupDeleteRequest struct {
	ApiKey      i.Apikey   `formam:"apiKey"`
	Rescource   i.Resource `formam:"resource"`
	Service     string     `formam:"service"`
	ServicePath string     `formam:"servicePath"`
}

func ConfigGroups(repo fr.FiwareRepo) chi.Router {
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
		serviceToConfigGroups := map[string][]i.ConfigGroup{}
		for _, v := range fss {
			sgs, err := iota.ListConfigGroups(*v)

			if err != nil {
				log.Err(err).Msg("Could not get fiware services")
				http.Error(w, "Could not get fiware services", http.StatusInternalServerError)
				return
			}
			if sgs.Count != 0 {
				serviceToConfigGroups[v.Service] = sgs.Services
			}
		}

		templates.Prepare(r, configgroup.FiwareServices(serviceToConfigGroups, "")).Render(r.Context(), w)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Bad Request"))
			components.Error(err).Render(r.Context(), w)
			return
		}

		var sg i.ConfigGroup
		var decoder = formam.NewDecoder(&formam.DecoderOptions{TagName: "schema"})
		err = decoder.Decode(r.PostForm, &sg)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(400)
			components.Error(err).Render(r.Context(), w)
			return
		}

		iota := i.IoTA{Host: globals.Conf.IoTAHost, Port: globals.Conf.IoTAPort}

		err = iota.CreateConfigGroup(i.FiwareService{Service: sg.Service, ServicePath: sg.ServicePath}, sg)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(500)
			components.Error(err).Render(r.Context(), w)
			return
		}
		configgroup.ConfigGroup(sg).Render(r.Context(), w)
		w.WriteHeader(http.StatusOK)

	})
	r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			components.Error(err).Render(r.Context(), w)
		}

		var req CofigGroupDeleteRequest
		var decoder = formam.NewDecoder(&formam.DecoderOptions{TagName: "formam"})
		err = decoder.Decode(r.Form, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			components.Error(err).Render(r.Context(), w)
		}
		iota := i.IoTA{Host: globals.Conf.IoTAHost, Port: globals.Conf.IoTAPort}
		err = iota.DeleteConfigGroup(i.FiwareService{Service: req.Service, ServicePath: req.ServicePath}, req.Rescource, req.ApiKey)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			components.Error(err).Render(r.Context(), w)
		}
		w.WriteHeader(http.StatusOK)
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
func AddConfigGroups(repo fr.FiwareRepo) chi.Router {
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
		templates.Prepare(r, configgroup.AddConfigGroupForm(string(encodedBytes))).Render(r.Context(), w)
	})
	return r
}
