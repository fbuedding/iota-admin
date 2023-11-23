package iotagentsdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/niemeyer/golang/src/pkg/container/vector"
	"github.com/rs/zerolog/log"
)

const (
	urlService = urlBase + "/iot/services"
)

func (e *MissingFields) Error() string {
	return fmt.Sprintf("Error %s: %s", e.Message, e.Fields)
}

func (sg ServiceGroup) Validate() error {

	mF := &MissingFields{make(vector.StringVector, 0), "Missing fields"}
	if sg.Apikey == "" {
		mF.Fields.Push("Apikey")
	}

	if mF.Fields.Len() == 0 {
		return nil
	} else {
		return mF
	}
}

type RespReadServiceGroup struct {
	Count    int            `json:"count"`
	Services []ServiceGroup `json:"services"`
}

type ReqCreateServiceGroup struct {
	Services []ServiceGroup `json:"services"`
}

func (i IoTA) ReadServiceGroup(fs FiwareService, r Resource, a Apikey) (*RespReadServiceGroup, error) {
	url := urlService + fmt.Sprintf("?resource=%s&apikey=%s", r, a)

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf(url, i.Host, i.Port), nil)

	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		json.Unmarshal(resData, &apiError)
		return nil, apiError
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}

	var respReadServiceGroup RespReadServiceGroup
	json.Unmarshal(responseData, &respReadServiceGroup)
	return &respReadServiceGroup, nil
}

func (i IoTA) ListServiceGroups(fs FiwareService) (*RespReadServiceGroup, error) {
	url := urlService

	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf(url, i.Host, i.Port), nil)

	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		json.Unmarshal(resData, &apiError)
		return nil, apiError
	}

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}

	var respReadServiceGroup RespReadServiceGroup
	json.Unmarshal(responseData, &respReadServiceGroup)
	return &respReadServiceGroup, nil
}

func (i IoTA) ServiceGroupExists(fs FiwareService, r Resource, a Apikey) bool {
	tmp, err := i.ReadServiceGroup(fs, r, a)
	if err != nil {
		return false
	}
	return tmp.Count > 0
}

func (i IoTA) CreateServiceGroup(fs FiwareService, sg ServiceGroup) error {
	sgs := [1]ServiceGroup{sg}
	return i.CreateServiceGroups(fs, sgs[:])
}

func (i IoTA) CreateServiceGroups(fs FiwareService, sgs []ServiceGroup) error {
	for _, sg := range sgs {
		err := sg.Validate()
		if err != nil {
			return err
		}
	}
	reqCreateServiceGroup := ReqCreateServiceGroup{}
	reqCreateServiceGroup.Services = sgs[:]
	method := "POST"

	payload, err := json.Marshal(reqCreateServiceGroup)
	if err != nil {
		log.Panic().Err(err).Msg("Could not Marshal struct")
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf(urlService, i.Host, i.Port), bytes.NewBuffer(payload))

	if err != nil {
		return fmt.Errorf("Error while creating Request %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while requesting resource %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		json.Unmarshal(resData, &apiError)
		return apiError
	}

	return nil
}

func (i IoTA) UpdateServiceGroup(fs FiwareService, r Resource, a Apikey, sg ServiceGroup) error {
	err := sg.Validate()
	if err != nil {
		return err
	}
	url := urlService + fmt.Sprintf("?resource=%s&apikey=%s", r, a)
	method := "PUT"

	payload, err := json.Marshal(sg)
	if err != nil {
		log.Panic().Err(err).Msg("Could not Marshal struct")
	}
	if string(payload) == "{}" {
		return nil
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf(url, i.Host, i.Port), bytes.NewBuffer(payload))

	if err != nil {
		return fmt.Errorf("Error while creating Request %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while requesting resource %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		json.Unmarshal(resData, &apiError)
		return apiError
	}

	return nil
}
func (i IoTA) DeleteServiceGroup(fs FiwareService, r Resource, a Apikey) error {
	url := urlService + fmt.Sprintf("?resource=%s&apikey=%s", r, a)

	method := http.MethodDelete

	client := http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf(url, i.Host, i.Port), strings.NewReader(""))

	if err != nil {
		return fmt.Errorf("Error while creating Request %w", err)
	}

	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error while requesting resource %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("Error while eding response body %w", err)
		}
		var apiError ApiError
		json.Unmarshal(resData, &apiError)
		return apiError
	}

	return nil
}

func (i IoTA) UpsertServiceGroup(fs FiwareService, sg ServiceGroup) {
	exists := i.ServiceGroupExists(fs, sg.Resource, sg.Apikey)
	if !exists {
		log.Debug().Msg("Creating service group...")
		err := i.CreateServiceGroup(fs, sg)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not create service group")
		}
	} else {
		log.Debug().Msg("Update service group...")
		err := i.UpdateServiceGroup(fs, sg.Resource, sg.Apikey, sg)
		if err != nil {
			log.Fatal().Err(err).Msg("Could not update service group")
		}
	}
}

func (i IoTA) CreateServiceGroupWSE(fs FiwareService, sg *ServiceGroup) error {
	if sg == nil {
		return errors.New("Service group reference cannot be nil")
	}

	err := i.CreateServiceGroup(fs, *sg)
	if err != nil {
    return err
	}

  sgTmp, err := i.ReadServiceGroup(fs, sg.Resource, sg.Apikey)
	if err != nil {
    return err
	}

  if sgTmp.Count == 0 {
    return errors.New("No service group created")
  }
  *sg = *&sgTmp.Services[0]

	return nil
}
