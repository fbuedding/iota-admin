package iotagentsdk

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	urlService = urlBase + "/iot/services"
)

type ReqCreateService struct {
	Apikey     string `json:"apikey"`
	Cbroker    string `json:"cbroker"`
	EntityType string `json:"entity_type"`
	Resource   string `json:"resource"`
}

func (i IoTA) GetService() (*RespHealthcheck, error) {

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
