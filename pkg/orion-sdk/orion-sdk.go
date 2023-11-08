package orionsdk

import (
	"encoding/json"
	"reflect"
	"strings"
)

type ID string

type Entity struct {
	ID    ID                   `json:"id"`
	Type  string               `json:"type"`
	Attrs map[string]Attribute `json:"-"`
}

type Attribute struct {
	Type     string              `json:"type"`
	Value    interface{}         `json:"value"`
	Metadata map[string]Metadata `json:"metadata"`
}

type Metadata struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (e Entity) MarshalJSON() ([]byte, error) {
  // quite simple, just put all into a map build a json string
	data := make(map[string]interface{})
	for k, v := range e.Attrs {
		data[k] = v
	}

	val := reflect.ValueOf(e)
	typ := reflect.TypeOf(e)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		fieldv := val.Field(i)
		jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonTag != "" && jsonTag != "-" {
			data[jsonTag] = fieldv.Interface()
		}
	}
	return json.Marshal(data)
}

func (e *Entity) UnmarshalJSON(b []byte) error {
  // less simple, first put Unmarshal all to entity struct. 
	type _tmp Entity
	e2 := _tmp{}
	err := json.Unmarshal(b, &e2) // this will set everything except the Attrs field

	if err != nil {
		return err
	}
  // now put all fields from the json bytes into an map of json.RawMessage
	var objmap map[string]json.RawMessage
	err = json.Unmarshal(b, &objmap)
	if err != nil {
		return err
	}
  // make the attribute map
	e2.Attrs = make(map[string]Attribute)
	typ := reflect.TypeOf(e2)
	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		jsonTag := strings.Split(field.Tag.Get("json"), ",")[0]
		if jsonTag != "" && jsonTag != "-" {
			delete(objmap, jsonTag) // delete the known Struct fields, thus everythin with an empty json tag or a "-"
		}
	}
  // now just put them into the Attrs map
	for k := range objmap {
		tmp := Attribute{}
		err = json.Unmarshal(objmap[k], &tmp)
		if err != nil {
			return err
		}
		e2.Attrs[k] = tmp
	}
  
	*e = Entity(e2)

	return nil
}
