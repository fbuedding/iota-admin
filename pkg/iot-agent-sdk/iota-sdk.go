package iotagentsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	urlBase        = "http://%v:%d"
	urlHealthcheck = urlBase + "/iot/about"
)

type IoTA struct {
	Host string
	Port int
}

type FiwareService struct {
	Service     string
	ServicePath string
}

type RespHealthcheck struct {
	LibVersion string `json:"libVersion"`
	Port       string `json:"port"`
	BaseRoot   string `json:"baseRoot"`
	Version    string `json:"version"`
}
type ApiError struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (e ApiError) Error() string {
  return fmt.Sprintf("%s: %s", e.Name, e.Message)
}

func (i IoTA) Healthcheck() (*RespHealthcheck, error) {
	response, err := http.Get(fmt.Sprintf(urlHealthcheck, i.Host, i.Port))
	if err != nil {
		return nil, fmt.Errorf("Error while Healthcheck: %w", err)
	}
	defer response.Body.Close()

	responseData, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, fmt.Errorf("Error while Healthcheck: %w", err)
	}
	var respHealth RespHealthcheck
	json.Unmarshal(responseData, &respHealth)
	return &respHealth, nil
}
