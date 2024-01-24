package routes

import (
	"encoding/json"
	"net/http"
	"slices"

	"github.com/fbuedding/iota-admin/internal/globals"
	"github.com/fbuedding/iota-admin/internal/helpers"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/fbuedding/iota-admin/web/templates"
	"github.com/fbuedding/iota-admin/web/templates/components"
	"github.com/fbuedding/iota-admin/web/templates/fiware/iotAgent/devices"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

type DevicDeleteRequest struct {
	ApiKey      i.Apikey   `formam:"apiKey"`
	Rescource   i.Resource `formam:"resource"`
	Service     string     `formam:"service"`
	ServicePath string     `formam:"servicePath"`
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
		iota := i.IoTA{Host: globals.Conf.IoTAHost, Port: globals.Conf.IoTAPort}
		fss := services.ToFiwareServices()
		log.Debug().Any("Fiware services", fss).Send()
		serviceToDevices := map[string][]i.Device{}
		for _, v := range fss {
			ds, err := iota.ListDevices(*v)

			if err != nil {
				log.Err(err).Msg("Could not get fiware services")
				http.Error(w, "Could not get fiware services", http.StatusInternalServerError)
				return
			}
			if ds.Count != 0 {
				serviceToDevices[v.Service] = ds.Devices
			}
		}

		log.Debug().Any("Devices", serviceToDevices).Send()

		templates.Prepare(r, devices.FiwareServices(serviceToDevices, "")).Render(r.Context(), w)
	})
	r.Post("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte("Bad Request"))
			components.Error(err).Render(r.Context(), w)
			return
		}

		var d i.Device
		err = helpers.Decode(r.PostForm, &d)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(400)
			components.Error(err).Render(r.Context(), w)
			return
		}

		iota := i.IoTA{Host: globals.Conf.IoTAHost, Port: globals.Conf.IoTAPort}

		err = iota.CreateDevice(i.FiwareService{Service: d.Service, ServicePath: d.ServicePath}, d)
		if err != nil {
			log.Error().Err(err).Send()
			w.WriteHeader(500)
			components.Error(err).Render(r.Context(), w)
			return
		}
		devices.Device(d).Render(r.Context(), w)
		w.WriteHeader(http.StatusOK)

	})
	r.Delete("/", func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			components.Error(err).Render(r.Context(), w)
		}

		var req CofigGroupDeleteRequest
		err = helpers.Decode(r.Form, &req)
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
		templates.Prepare(r, devices.AddDeviceForm(string(encodedBytes))).Render(r.Context(), w)
	})
	return r
}
