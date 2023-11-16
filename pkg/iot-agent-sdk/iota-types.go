package iotagentsdk

import "github.com/niemeyer/golang/src/pkg/container/vector"

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
	ObjectID   string              `json:"object_id,omitempty"`
	Name       string              `json:"name"`
	Type       string              `json:"type"`
	Metadata   map[string]Metadata `json:"metadata,omitempty"`
	Expression string              `json:"expression,omitempty"`
	SkipValue  string              `json:"skipValue,omitempty"`
	EntityName string              `json:"entity_name,omitempty"`
	EntityType string              `json:"entity_type,omitempty"`
}

type LazyAttribute struct {
	ObjectID string              `json:"object_id,omitempty"`
	Name     string              `json:"name"`
	Type     string              `json:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty"`
}

type StaticAttribute struct {
	ObjectID string              `json:"object_id,omitempty"`
	Name     string              `json:"name"`
	Type     string              `json:"type"`
	Metadata map[string]Metadata `json:"metadata,omitempty"`
}

type Command struct {
	ObjectID    string              `json:"object_id,omitempty"`
	Name        string              `json:"name"`
	Type        string              `json:"type"`
	Expression  string              `json:"expression,omitempty"`
	PayloadType string              `json:"payloadType,omitempty"`
	ContentType string              `json:"contentType,omitempty"`
	Metadata    map[string]Metadata `json:"metadata,omitempty"`
}

type Metadata struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Apikey string
type Resource string

// see https://iotagent-node-lib.readthedocs.io/en/latest/api.html#service-group-datamodel
type ServiceGroup struct {
	Service                      string            `json:"service,omitempty"`
	ServicePath                  string            `json:"subservice,omitempty"`
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

type DeciveId string

type Device struct {
	Id                 DeciveId          `json:"device_id"`
	Service            string            `json:"service,omitempty"`
	ServicePath        string            `json:"service_path,omitempty"`
	EntityName         string            `json:"enitity_name"`
	Type               string            `json:"type,omitempty"`
	EntityType         string            `json:"entity_type,omitempty"` //both type and EntityType work and are used both in documentation
	Timezone           string            `json:"timezon,omitempty"`
	Timestamp          *bool             `json:"timestamp,omitempty"`
	Apikey             Apikey            `json:"apikey"`
	Endpoint           string            `json:"endpoint,omitempty"`
	Protocol           string            `json:"protocol,omitempty"`
	Transport          string            `json:"transport,omitempty"`
	Attributes         []Attribute       `json:"attributes,omitempty"`
	Commands           []Command         `json:"commands,omitempty"`
	Lazy               []LazyAttribute   `json:"lazy,omitempty"`
	StaticAttributes   []StaticAttribute `json:"static_attributes,omitempty"`
	InternalAttributes []interface{}     `json:"internal_attributes,omitempty"`
	ExplicitAttrs      string            `json:"explicitAttrs,omitempty"`
	NgsiVersion        string            `json:"ngsiVersion,omitempty"`
}

type MissingFields struct {
	Fields  vector.StringVector
	Message string
}
