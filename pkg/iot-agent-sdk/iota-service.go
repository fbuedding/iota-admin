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

func (sg ConfigGroup) Validate() error {

	mF := &MissingFields{make(vector.StringVector, 0), "Missing fields"}
	if sg.Apikey == "" {
		mF.Fields.Push("Apikey")
	}
	if sg.Resource == "" {
		mF.Fields.Push("Resource")
	}

	if mF.Fields.Len() == 0 {
		return nil
	} else {
		return mF
	}
}

type RespReadConfigGroup struct {
	Count    int           `json:"count"`
	Services []ConfigGroup `json:"services"`
}

type ReqCreateConfigGroup struct {
	Services []ConfigGroup `json:"services"`
}

func (i IoTA) ReadConfigGroup(fs FiwareService, r Resource, a Apikey) (*RespReadConfigGroup, error) {
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

	var respReadConfigGroup RespReadConfigGroup
	json.Unmarshal(responseData, &respReadConfigGroup)
	return &respReadConfigGroup, nil
}

func (i IoTA) ListConfigGroups(fs FiwareService) (*RespReadConfigGroup, error) {
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

	var respReadConfigGroup RespReadConfigGroup
	json.Unmarshal(responseData, &respReadConfigGroup)
	return &respReadConfigGroup, nil
}

func (i IoTA) ConfigGroupExists(fs FiwareService, r Resource, a Apikey) bool {
	tmp, err := i.ReadConfigGroup(fs, r, a)
	if err != nil {
		return false
	}
	return tmp.Count > 0
}

func (i IoTA) CreateConfigGroup(fs FiwareService, sg ConfigGroup) error {
	sgs := [1]ConfigGroup{sg}
	return i.CreateConfigGroups(fs, sgs[:])
}

func (i IoTA) CreateConfigGroups(fs FiwareService, sgs []ConfigGroup) error {
	for _, sg := range sgs {
		err := sg.Validate()
		if err != nil {
			return err
		}
	}
	reqCreateConfigGroup := ReqCreateConfigGroup{}
	reqCreateConfigGroup.Services = sgs[:]
	method := "POST"

	payload, err := json.Marshal(reqCreateConfigGroup)
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

func (i IoTA) UpdateConfigGroup(fs FiwareService, r Resource, a Apikey, sg ConfigGroup) error {
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
func (i IoTA) DeleteConfigGroup(fs FiwareService, r Resource, a Apikey) error {
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

func (i IoTA) UpsertConfigGroup(fs FiwareService, sg ConfigGroup) error {
	exists := i.ConfigGroupExists(fs, sg.Resource, sg.Apikey)
	if !exists {
		log.Debug().Msg("Creating service group...")
		err := i.CreateConfigGroup(fs, sg)
		if err != nil {
			return err
		}
	} else {
		log.Debug().Msg("Update service group...")
		err := i.UpdateConfigGroup(fs, sg.Resource, sg.Apikey, sg)
		if err != nil {
			return err
		}
	}
	return nil
}

func (i IoTA) CreateConfigGroupWSE(fs FiwareService, sg *ConfigGroup) error {
	if sg == nil {
		return errors.New("Service group reference cannot be nil")
	}

	err := i.CreateConfigGroup(fs, *sg)
	if err != nil {
		return err
	}

	sgTmp, err := i.ReadConfigGroup(fs, sg.Resource, sg.Apikey)
	if err != nil {
		return err
	}

	if sgTmp.Count == 0 {
		return errors.New("No service group created")
	}
	*sg = *&sgTmp.Services[0]

	return nil
}
