package orionsdk_test

import (
	"testing"

	orionsdk "github.com/fbuedding/iota-admin/pkg/orion-sdk"
)

func TestMain(m *testing.M) {
	m.Run()
}

func TestMarshalEntity(t *testing.T) {
	entity := orionsdk.Entity{
		ID:   "Test",
		Type: "Test",
		Attrs: map[string]orionsdk.Attribute{"test": {
			Type:  "Number",
			Value: 2,
			Metadata: map[string]orionsdk.Metadata{"test": {
				Type:  "Bli",
				Value: "bla",
			}},
		}},
	}
	b, err := entity.MarshalJSON()
	if err != nil {
		t.Fail()
	}
	t.Log(string(b))
}
