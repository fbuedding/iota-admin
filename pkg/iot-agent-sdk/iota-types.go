package iotagentsdk

import (
	"github.com/niemeyer/golang/src/pkg/container/vector"
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

// these are nearly the same, but for typesafety differnt structs
type Attribute struct {
	ObjectID   string              `json:"object_id,omitempty" schema:"object_id"`
	Name       string              `json:"name" schema:"name"`
	Type       string              `json:"type" schema:"type"`
	Metadata   map[string]Metadata `json:"metadata,omitempty" schema:"metadata"`
	Expression string              `json:"expression,omitempty" schema:"expression"`
	SkipValue  string              `json:"skipValue,omitempty" schema:"skipValue"`
	EntityName string              `json:"entity_name,omitempty" schema:"entity_name"`
	EntityType string              `json:"entity_type,omitempty" schema:"entity_type"`
}

type LazyAttribute struct {
	ObjectID string              `json:"object_id,omitempty" schema:"object_id"`
	Name     string              `json:"name" schema:"name"`
	Type     string              `json:"type" schema:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty" schema:"metadata"`
}

type StaticAttribute struct {
	ObjectID string              `json:"object_id,omitempty"`
	Name     string              `json:"name" schema:"name"`
	Type     string              `json:"type" schema:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty"`
}

type Command struct {
	ObjectID    string              `json:"object_id,omitempty" schema:"object_id"`
	Name        string              `json:"name" schema:"name"`
	Type        string              `json:"type" schema:"type"`
	Expression  string              `json:"expression,omitempty" schema:"expression"`
	PayloadType string              `json:"payloadType,omitempty" schema:"payloadType"`
	ContentType string              `json:"contentType,omitempty" schema:"contentType"`
	Metadata    map[string]Metadata `json:"metadata,omitempty" schema:"metadata"`
}

// Extra type for x-form-data
type Metadata struct {
	Type  string `json:"type" schema:"type"`
	Value string `json:"value" schema:"value"`
}

type Apikey string
type Resource string

// see https://iotagent-node-lib.readthedocs.io/en/latest/api.html#service-group-datamodel
type ServiceGroup struct {
	Service                      string            `json:"service,omitempty" schema:"service"`
	ServicePath                  string            `json:"subservice,omitempty" schema:"subservice"`
	Resource                     Resource          `json:"resource" schema:"resource"`
	Apikey                       Apikey            `json:"apikey" schema:"apikey"`
	Timestamp                    *bool             `json:"timestamp,omitempty" schema:"timestamp"`
	EntityType                   string            `json:"entity_type,omitempty" schema:"entity_type"`
	Trust                        string            `json:"trust,omitempty" schema:"trust"`
	CbHost                       string            `json:"cbHost,omitempty" schema:"cbHost"`
	Lazy                         []LazyAttribute   `json:"lazy,omitempty" schema:"lazy"`
	Commands                     []Command         `json:"commands,omitempty" schema:"commands"`
	Attributes                   []Attribute       `json:"attributes,omitempty" schema:"attributes"`
	StaticAttributes             []StaticAttribute `json:"static_attributes,omitempty" schema:"static_attributes"`
	InternalAttributes           []interface{}     `json:"internal_attributes,omitempty" schema:"internal_attributes"`
	ExplicitAttrs                string            `json:"explicitAttrs,omitempty" schema:"explicitAttrs"`
	EntityNameExp                string            `json:"entityNameExp,omitempty" schema:"entityNameExp"`
	NgsiVersion                  string            `json:"ngsiVersion,omitempty" schema:"ngsiVersion"`
	DefaultEntityNameConjunction string            `json:"defaultEntityNameConjunction,omitempty" schema:"defaultEntityNameConjunction"`
	Autoprovision                bool              `json:"autoprovision,omitempty" schema:"autoprovision"`
}

type DeciveId string

type Device struct {
	Id                 DeciveId          `json:"device_id,omitempty" schema:"device_id"`
	Service            string            `json:"service,omitempty" schema:"service"`
	ServicePath        string            `json:"service_path,omitempty" schema:"service_path"`
	EntityName         string            `json:"entity_name,omitempty" schema:"entity_name"`
	EntityType         string            `json:"entity_type,omitempty" schema:"entity_type"`
	Timezone           string            `json:"timezon,omitempty" schema:"timezon"`
	Timestamp          *bool             `json:"timestamp,omitempty" schema:"timestamp"`
	Apikey             Apikey            `json:"apikey,omitempty" schema:"apikey"`
	Endpoint           string            `json:"endpoint,omitempty" schema:"endpoint"`
	Protocol           string            `json:"protocol,omitempty" schema:"protocol"`
	Transport          string            `json:"transport,omitempty" schema:"transport"`
	Attributes         []Attribute       `json:"attributes,omitempty" schema:"attributes"`
	Commands           []Command         `json:"commands,omitempty" schema:"commands"`
	Lazy               []LazyAttribute   `json:"lazy,omitempty" schema:"lazy"`
	StaticAttributes   []StaticAttribute `json:"static_attributes,omitempty" schema:"static_attributes"`
	InternalAttributes []interface{}     `json:"internal_attributes,omitempty" schema:"internal_attributes"`
	ExplicitAttrs      string            `json:"explicitAttrs,omitempty" schema:"explicitAttrs"`
	NgsiVersion        string            `json:"ngsiVersion,omitempty" schema:"ngsiVersion"`
}

type MissingFields struct {
	Fields  vector.StringVector
	Message string
}
