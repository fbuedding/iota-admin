package iotagentsdk

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	u "net/url"

	"github.com/niemeyer/golang/src/pkg/container/vector"
)

const (
	urlDevice = urlBase + "/iot/services"
)

type reqCreateDevice struct {
	Devices []Device `json:"devices"`
}

func (d Device) Validate() error {

	mF := &MissingFields{make(vector.StringVector, 0), "Missing fields"}
	if d.EntityType == "" && d.Type == "" {
		mF.Fields.Push("EntityType|Type")
	}

	if mF.Fields.Len() == 0 {
		return nil
	} else {
		return mF
	}
}

func (i IoTA) ReadDevice(fs FiwareService, id DeciveId) (*Device, error) {
	url, err := u.JoinPath(fmt.Sprintf(urlDevice, i.Host, i.Port), u.PathEscape(string(id)))

	if err != nil {
		return nil, err
	}
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}
	req.Header.Add("fiware-service", fs.Service)
	req.Header.Add("fiware-servicepath", fs.ServicePath)

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

	var device Device
	json.Unmarshal(responseData, &device)
	return &device, nil
}

func (i IoTA) DeviceExists(fs FiwareService, id DeciveId) bool {
	_, err := i.ReadDevice(fs, id)
	if err != nil {
		return false
	}
	return true
}

func (i IoTA) CreateDevices(fs FiwareService, ds []Device) error {

	for _, sg := range ds {
		err := sg.Validate()
		if err != nil {
			return err
		}
	}
	rcd := reqCreateDevice{}
	rcd.Devices = ds[:]
	method := "POST"

	payload, err := json.Marshal(rcd)
	if err != nil {
		fmt.Println("Could not Marshal struct")
		panic(1)
	}
	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf(urlDevice, i.Host, i.Port), bytes.NewBuffer(payload))

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
func (i IoTA) CreateDevice(fs FiwareService, d Device) error {
	ds := [1]Device{d}
	return i.CreateDevices(fs, ds[:])
}
