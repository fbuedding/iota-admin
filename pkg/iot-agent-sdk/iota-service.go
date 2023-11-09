package iotagentsdk

import (
	"encoding/json"
	"fmt"
	"github.com/niemeyer/golang/src/pkg/container/vector"
	"io"
	"net/http"
)

const (
	urlService = urlBase + "/iot/services"
)

// see https://iotagent-node-lib.readthedocs.io/en/latest/api.html#service-group-datamodel
type ServiceGroup struct {
	Resource                     Resource          `json:"resource"`
	Apikey                       Apikey            `json:"apikey"`
	Timestamp                    *bool             `json:"timestamp,omitempty"`
	Type                         string            `json:"type,omitempty"`
	EntityType                   string            `json:"entity_type,omitempty"` //both type and EntityType work and are used both in documentation
	Trust                        string            `json:"trust,omitempty"`
	CbHost                       string            `json:"cbHost,omitempty"`
	Lazy                         []LazyAttribute   `json:"lazy,omitempty"`
	Commands                     []Command         `json:"commands,omitempty"`
	Attributes                   []Attribute       `json:"attributes,omitempty"`
	StaticAttributes             []StaticAttribute `json:"static_attributes,omitempty"`
	InternalAttributes           []interface{}     `json:"internal_attributes,omitempty"`
	ExplicitAttrs                string            `json:"explicitAttrs,omitempty"`
	EntityNameExp                string            `json:"entityNameExp,omitempty"`
	NgsiVersion                  string            `json:"ngsiVersion,omitempty"`
	DefaultEntityNameConjunction string            `json:"defaultEntityNameConjunction,omitempty"`
	Autoprovision                bool              `json:"autoprovision,omitempty"`
}

type Apikey string
type Resource string

// these are all the same, but for typesafety differnt structs

type Command struct {
	Name     string              `json:"name"`
	Type     string              `json:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty"`
}

type Attribute struct {
	Name     string              `json:"name"`
	Type     string              `json:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty"`
}
type LazyAttribute struct {
	Name     string              `json:"name"`
	Type     string              `json:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty"`
}
type StaticAttribute struct {
	Name     string              `json:"name"`
	Type     string              `json:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type MissingFields struct {
	Fields  vector.StringVector
	Message string
}

func (e *MissingFields) Error() string {
	return fmt.Sprintf("Error %s: %s", e.Message, e.Fields)
}

func (sg ServiceGroup) Validate() error {

	mF := &MissingFields{make(vector.StringVector, 0), "Missing fields"}

	if sg.Resource == "" {
		mF.Fields.Push("Resource")
	}
	if sg.Apikey == "" {
		mF.Fields.Push("Apikey")
	}
	if sg.EntityType == "" && sg.Type == "" {
		mF.Fields.Push("EntityType|Type")
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

func (i IoTA) ReadServiceGroup(fs FiwareService, r Resource, a Apikey) (*RespReadServiceGroup, error) {
	url := urlService + fmt.Sprintf("?r=%s&a=%s", r, a)

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

	responseData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("Error while getting service: %w", err)
	}

	var respReadServiceGroup RespReadServiceGroup
	json.Unmarshal(responseData, &respReadServiceGroup)
	return &respReadServiceGroup, nil
}

func (i IoTA) ServiceGroupExists(fs FiwareService, r Resource, a Apikey) (bool, error) {
	tmp, err := i.ReadServiceGroup(fs, r, a)
	if err != nil {
		return false, err
	}
  return tmp.Count > 0, nil
}


func (i IoTA) CreateServiceGroup(fs FiwareService, sg ServiceGroup) error {
    err := sg.Validate()
    if err != nil {
      return err
    }
    return nil
}
