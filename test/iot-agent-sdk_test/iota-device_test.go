package iotagentsdktest_test

import (
	"log"
	"testing"

	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
)

var (
	iota i.IoTA
	fs   i.FiwareService
	d    i.Device
)

const (
	deviceId = i.DeciveId("test_device")
)

func init() {
	iota = i.IoTA{Host: "localhost", Port: 4061}
	fs = i.FiwareService{Service: "testing", ServicePath: "/"}
	d = i.Device{Id: deviceId}
  err := iota.CreateDevice(fs, d) 
  if err != nil {
    log.Fatal(err.Error())
  }
}
func TestUpdateDevice(t *testing.T) {
  dtmp := d
  dtmp.EntityName = "test-entity-name-updated"
  dtmp.EntityType = "test"
  err := iota.UpdateDevice(fs, dtmp)  
  if err != nil {
    t.Error(err)
  }
}

func TestDeleteDevice(t *testing.T) {
  err := iota.DeleteDevice(fs, d.Id)
  if err != nil {
    t.Error(err)
  }
}
