package iotagentsdktest_test

import (
	"testing"

	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/rs/zerolog/log"
)

var (
	iota i.IoTA
	fs   i.FiwareService
	d    i.Device
)

const (
	deviceId          = i.DeciveId("test_device")
	entityName        = "TestEntityName"
	updatedEntityName = "TestEntityNameUpdated"
)

func init() {
	iota = i.IoTA{Host: "localhost", Port: 4061}
	fs = i.FiwareService{Service: "testing", ServicePath: "/"}
	d = i.Device{Id: deviceId, EntityName: entityName}
	err := iota.CreateDevice(fs, d)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create device")
	}
}

func TestReadDevice(t *testing.T) {
	respD, err := iota.ReadDevice(fs, deviceId)
	if err != nil {
		t.Error(err)
	}
	if respD.EntityName != entityName {
		t.Fail()
	}
}
func TestListDevice(t *testing.T) {
	respD, err := iota.ListDevices(fs)
	if err != nil {
		t.Error(err)
	}
	if respD.Count != 1 {
		t.Fail()
	}

	if respD.Devices[0].EntityName != entityName {
		t.Fail()
	}
}

func TestUpdateDevice(t *testing.T) {
	dtmp := d
	dtmp.EntityName = updatedEntityName
	dtmp.EntityType = "test"
	dtmp.Transport = "MQTT"
	err := iota.UpdateDevice(fs, dtmp)
	if err != nil {
		t.Error(err)
	}
	dUpdated, _ := iota.ReadDevice(fs, d.Id)
	if dUpdated.EntityName != updatedEntityName {
		t.Fail()
	}
	dtmp1 := i.Device{Id: deviceId}
  
	err = iota.UpdateDevice(fs, dtmp1)

	if err != nil {
    t.Log("Device shouldn't updatet empty body")
		t.Error(err)
	}
}

func TestDeleteDevice(t *testing.T) {
	err := iota.DeleteDevice(fs, d.Id)
	if err != nil {
		t.Error(err)
	}
}
