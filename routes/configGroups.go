package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"github.com/fbuedding/iota-admin/internal/globals"
	"github.com/fbuedding/iota-admin/internal/helpers"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/fbuedding/iota-admin/web/templates"
	configgroup "github.com/fbuedding/iota-admin/web/templates/fiware/iotAgent/configGroup"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type ConfigGroupDeleteRequest struct {
	IoTAgentId  string     `formam:"iotAgentId"`
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
		iotas, err := repo.ListIotas()
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
		}
		fss := services.ToFiwareServices()
		iotAgentToServiceToConfigGroups := configgroup.IoTAToFiwareServiceToConfigGroupsWithIoTAId{}
		for _, iota := range iotas {
			if _, err := iota.ToIoTA().Healthcheck(); err != nil {
				continue
			}
			iotAgentToServiceToConfigGroups[iota.Alias] = configgroup.FiwareServiceToConfigGroupsWithIoTAId{}
			for _, v := range fss {
				sgs, err := iota.ToIoTA().ListConfigGroups(*v)
				if err != nil {
					log.Err(err).Msg("Could not get fiware services")
				}
				if sgs.Count != 0 {
					iotAgentToServiceToConfigGroups[iota.Alias][v.Service] = configgroup.ConfigGroupsWithIoTAId{
						IoTAId:       iota.Id,
						ConfigGroups: sgs.Services,
					}
				}
			}
		}

		templates.Prepare(r, configgroup.IoTAgents(iotAgentToServiceToConfigGroups)).Render(r.Context(), w)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusBadRequest)
			return
		}
		if !r.PostForm.Has("iotAgent") {
			templates.HandleError(r.Context(), w, fmt.Errorf("No Iot-Agent provided"), http.StatusUnprocessableEntity)
			return
		}

		iotaId := r.PostForm.Get("iotAgent")
		r.PostForm.Del("iotAgent")
		iota, err := repo.GetIota(iotaId)
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			log.Error().Err(err).Msgf("Could not get IoT-Agent for id: %s", iotaId)
			return
		}
		var sg i.ConfigGroup
		err = helpers.Decode(r.PostForm, &sg)
		if err != nil {
			log.Error().Err(err).Send()
			templates.HandleError(r.Context(), w, err, http.StatusBadRequest)
			return
		}

		err = iota.CreateConfigGroup(i.FiwareService{Service: sg.Service, ServicePath: sg.ServicePath}, sg)
		if err != nil {
			log.Error().Err(err).Send()
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			return
		}
		configgroup.ConfigGroup(sg, iotaId).Render(r.Context(), w)
		w.WriteHeader(http.StatusOK)
	})
	r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusBadRequest)
		}

		var req ConfigGroupDeleteRequest
		err = helpers.Decode(r.Form, &req)
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusBadRequest)
		}
		iota, err := repo.GetIota(req.IoTAgentId)
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			log.Error().Err(err).Msgf("Could not get IoT-Agent for id: %s", req.IoTAgentId)
			return
		}
		err = iota.DeleteConfigGroup(i.FiwareService{Service: req.Service, ServicePath: req.ServicePath}, req.Rescource, req.ApiKey)
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
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

func AddConfigGroupForm(repo fr.FiwareRepo) chi.Router {
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
		iotAgents, err := repo.ListIotas()
		if err != nil {
			templates.HandleError(r.Context(), w, err, 500)
			log.Error().Err(err).Msgf("Could not get iot-agents")
		}
		encodedIotAgents, err := json.Marshal(iotAgents)
		if err != nil {
			templates.HandleError(r.Context(), w, err, 500)
			log.Error().Err(err).Msgf("Could not encoded iot-agents")
		}
		templates.Prepare(r, configgroup.AddConfigGroupForm(string(encodedBytes), string(encodedIotAgents))).Render(r.Context(), w)
	})
	return r
}
