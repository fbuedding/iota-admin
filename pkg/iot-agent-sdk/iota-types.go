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

	ObjectID   string              `json:"object_id,omitempty" formam:"object_id"`
	Name       string              `json:"name" formam:"name"`
	Type       string              `json:"type" formam:"type"`
	Expression string              `json:"expression,omitempty" formam:"expression"`
	SkipValue  string              `json:"skipValue,omitempty" formam:"skipValue"`
	EntityName string              `json:"entity_name,omitempty" formam:"entity_name"`
	EntityType string              `json:"entity_type,omitempty" formam:"entity_type"`
	Metadata   map[string]Metadata `json:"metadata,omitempty" formam:"metadata"`
}

type LazyAttribute struct {
	ObjectID string              `json:"object_id,omitempty" formam:"object_id"`
	Name     string              `json:"name" formam:"name"`
	Type     string              `json:"type" formam:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty" formam:"metadata"`
}

type StaticAttribute struct {
	ObjectID string              `json:"object_id,omitempty" formam:"object_id"`
	Name     string              `json:"name" formam:"name"`
	Type     string              `json:"type" formam:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty" formam:"metadata"`
}

type Command struct {
	ObjectID    string              `json:"object_id,omitempty" formam:"object_id"`
	Name        string              `json:"name" formam:"name"`
	Type        string              `json:"type" formam:"type"`
	Expression  string              `json:"expression,omitempty" formam:"expression"`
	PayloadType string              `json:"payloadType,omitempty" formam:"payloadType"`
	ContentType string              `json:"contentType,omitempty" formam:"contentType"`
	Metadata    map[string]Metadata `json:"metadata,omitempty" formam:"metadata"`
}

// Extra type for x-form-data
type Metadata struct {
	Type  string `json:"type" formam:"type"`
	Value string `json:"value" formam:"value"`
}

type Apikey string
type Resource string

// see https://iotagent-node-lib.readthedocs.io/en/latest/api.html#service-group-datamodel
type ConfigGroup struct {
	Service                      string            `json:"service,omitempty" formam:"service"`
	ServicePath                  string            `json:"subservice,omitempty" formam:"subservice"`
	Resource                     Resource          `json:"resource" formam:"resource"`
	Apikey                       Apikey            `json:"apikey" formam:"apikey"`
	Timestamp                    *bool             `json:"timestamp,omitempty" formam:"timestamp"`
	EntityType                   string            `json:"entity_type,omitempty" formam:"entity_type"`
	Trust                        string            `json:"trust,omitempty" formam:"trust"`
	CbHost                       string            `json:"cbHost,omitempty" formam:"cbHost"`
	Lazy                         []LazyAttribute   `json:"lazy,omitempty" formam:"lazy"`
	Commands                     []Command         `json:"commands,omitempty" formam:"commands"`
	Attributes                   []Attribute       `json:"attributes,omitempty" formam:"attributes"`
	StaticAttributes             []StaticAttribute `json:"static_attributes,omitempty" formam:"static_attributes"`
	InternalAttributes           []interface{}     `json:"internal_attributes,omitempty" formam:"internal_attributes"`
	ExplicitAttrs                string            `json:"explicitAttrs,omitempty" formam:"explicitAttrs"`
	EntityNameExp                string            `json:"entityNameExp,omitempty" formam:"entityNameExp"`
	NgsiVersion                  string            `json:"ngsiVersion,omitempty" formam:"ngsiVersion"`
	DefaultEntityNameConjunction string            `json:"defaultEntityNameConjunction,omitempty" formam:"defaultEntityNameConjunction"`
	Autoprovision                bool              `json:"autoprovision,omitempty" formam:"autoprovision"`
	PayloadType                  string            `json:"payloadType,omitempty" formam:"payloadType"`
	Transport                    string            `json:"transport,omitempty" formam:"transport"`
	Endpoint                     string            `json:"endpoint,omitempty" formam:"endpoint"`
}

type DeciveId string

type Device struct {
	Id                 DeciveId          `json:"device_id,omitempty" formam:"device_id"`
	Service            string            `json:"service,omitempty" formam:"service"`
	ServicePath        string            `json:"service_path,omitempty" formam:"service_path"`
	EntityName         string            `json:"entity_name,omitempty" formam:"entity_name"`
	EntityType         string            `json:"entity_type,omitempty" formam:"entity_type"`
	Timezone           string            `json:"timezon,omitempty" formam:"timezon"`
	Timestamp          *bool             `json:"timestamp,omitempty" formam:"timestamp"`
	Apikey             Apikey            `json:"apikey,omitempty" formam:"apikey"`
	Endpoint           string            `json:"endpoint,omitempty" formam:"endpoint"`
	Protocol           string            `json:"protocol,omitempty" formam:"protocol"`
	Transport          string            `json:"transport,omitempty" formam:"transport"`
	Attributes         []Attribute       `json:"attributes,omitempty" formam:"attributes"`
	Commands           []Command         `json:"commands,omitempty" formam:"commands"`
	Lazy               []LazyAttribute   `json:"lazy,omitempty" formam:"lazy"`
	StaticAttributes   []StaticAttribute `json:"static_attributes,omitempty" formam:"static_attributes"`
	InternalAttributes []interface{}     `json:"internal_attributes,omitempty" formam:"internal_attributes"`
	ExplicitAttrs      string            `json:"explicitAttrs,omitempty" formam:"explicitAttrs"`
	NgsiVersion        string            `json:"ngsiVersion,omitempty" formam:"ngsiVersion"`
	PayloadType        string            `json:"payloadType,omitempty" formam:"payloadType"`
}

type MissingFields struct {
	Fields  vector.StringVector
	Message string
}
