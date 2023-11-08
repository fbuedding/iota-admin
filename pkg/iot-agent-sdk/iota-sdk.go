package iotagentsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
  urlBase = "http://%v:%d"
  urlHealthcheck = urlBase + "/iot/about"
)
type IoTA struct {
	Host string
	Port int
}

type FiwareService struct {
  Service string
  ServicePath string
}

type RespHealthcheck struct {
	LibVersion string `json:"libVersion"`
	Port       string `json:"port"`
	BaseRoot   string `json:"baseRoot"`
	Version    string `json:"version"`
}

func (i IoTA) Healthcheck() (*RespHealthcheck, error) {
	response, err := http.Get(fmt.Sprintf(urlHealthcheck, i.Host, i.Port))
	if err != nil {
		return nil, fmt.Errorf("Error while Healthcheck: %w", err)
	}

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("Error while Healthcheck: %w", err)
	}
  var respHealth RespHealthcheck
  json.Unmarshal(responseData, &respHealth)
  return &respHealth, nil
}
