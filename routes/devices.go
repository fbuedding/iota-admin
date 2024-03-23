package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	i "github.com/fbuedding/fiware-iot-agent-sdk"
	"github.com/fbuedding/iota-admin/internal/globals"
	"github.com/fbuedding/iota-admin/internal/helpers"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/web/templates"
	"github.com/fbuedding/iota-admin/web/templates/fiware/iotAgent/devices"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type DevicDeleteRequest struct {
	DeviceId    i.DeciveId `formam:"deviceId"`
	Service     string     `formam:"service"`
	ServicePath string     `formam:"servicePath"`
	IoTAgentId  string     `formam:"iotAgentId"`
}

func Devices(repo fr.FiwareRepo) chi.Router {
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
		log.Debug().Any("Fiware services", fss).Send()
		iotAgentToServiceToDevices := devices.IoTAToFiwareServiceToDevicesWithIoTAId{}
		for _, iotaRow := range iotas {
			iota := iotaRow.ToIoTA()
			if _, err := iota.Healthcheck(); err != nil {
				continue
			}
			iotAgentToServiceToDevices[iotaRow.Alias] = devices.FiwareServiceToDevicesWithIoTAId{}
			for _, fiwareService := range fss {
				servicePaths, err := iota.GetAllServicePathsForService(fiwareService.Service)
				if err != nil {
					log.Err(err).Msg("Could not get Service-Paths for service")
				}
				for _, servicePath := range servicePaths {
					fiwareService.ServicePath = servicePath
					ds, err := iota.ListDevices(*fiwareService)
					if err != nil {
						log.Err(err).Msg("Could not get fiware services")
						http.Error(w, "Could not get fiware services", http.StatusInternalServerError)
						return
					}
					if ds.Count != 0 {
						iotAgentToServiceToDevices[iotaRow.Alias][fiwareService.Service] = devices.DevicesWithIoTAId{IoTAId: iotaRow.Id, Devices: ds.Devices}
					}
				}
			}
		}

		log.Debug().Any("Devices", iotAgentToServiceToDevices).Send()

		templates.Prepare(r, devices.IoTAgents(iotAgentToServiceToDevices)).Render(r.Context(), w)
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
		var d i.Device
		err = helpers.Decode(r.PostForm, &d)
		if err != nil {
			log.Error().Err(err).Send()
			templates.HandleError(r.Context(), w, err, http.StatusBadRequest)
			return
		}

		err = iota.CreateDevice(i.FiwareService{Service: d.Service, ServicePath: d.ServicePath}, d)
		if err != nil {
			log.Error().Err(err).Send()
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			return
		}
		devices.Device(d, iotaId).Render(r.Context(), w)
		w.WriteHeader(http.StatusOK)
	})
	r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			log.Err(err).Msg("Error while parsing form")
			templates.HandleError(r.Context(), w, err, http.StatusBadRequest)
		}

		var req DevicDeleteRequest
		err = helpers.Decode(r.Form, &req)
		if err != nil {
			log.Err(err).Msg("Error while decoding request")
			templates.HandleError(r.Context(), w, err, http.StatusBadRequest)
			return
		}

		iota, err := repo.GetIota(req.IoTAgentId)
		if err != nil {
			templates.HandleError(r.Context(), w, err, http.StatusInternalServerError)
			log.Error().Err(err).Msgf("Could not get IoT-Agent for id: %s", req.IoTAgentId)
			return
		}

		err = iota.DeleteDevice(i.FiwareService{
			Service:     req.Service,
			ServicePath: req.ServicePath,
		}, req.DeviceId)
		if err != nil {
			log.Err(err).Msg("Error while deleteing device")
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

func AddDeviceForm(repo fr.FiwareRepo) chi.Router {
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
		templates.Prepare(r, devices.AddDeviceForm(string(encodedBytes), string(encodedIotAgents))).Render(r.Context(), w)
	})
	return r
}
