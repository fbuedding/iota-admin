package iotagentsdktest_test

import (
	"testing"

	i "github.com/fbuedding/fiware-iot-agent-sdk"
	"github.com/rs/zerolog/log"
)

var (
	iota i.IoTA
	fs   i.FiwareService
	d    i.Device
	sg   i.ConfigGroup
)

const (
	deviceId          = i.DeciveId("test_device")
	entityName        = "TestEntityName"
	updatedEntityName = "TestEntityNameUpdated"
	service           = "testing"
	servicePath       = "/"
	resource          = i.Resource("/iot/d")
	apiKey            = "testKey"
)

func TestMain(m *testing.M) {
	iota = i.IoTA{Host: "localhost", Port: 4061}
	fs = i.FiwareService{Service: service, ServicePath: servicePath}
	d = i.Device{Id: deviceId, EntityName: entityName}
	sg = i.ConfigGroup{
		Service:       service,
		ServicePath:   servicePath,
		Resource:      resource,
		Apikey:        apiKey,
		Autoprovision: false,
	}
	_ = iota.DeleteDevice(fs, d.Id)
	err := iota.CreateDevice(fs, d)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create device for tests")
	}
	_ = iota.DeleteConfigGroup(fs, resource, apiKey)
	err = iota.CreateConfigGroup(fs, sg)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create service group for tests")
	}
	m.Run()
	teardown()
}

func teardown() {
	err := iota.DeleteDevice(fs, d.Id)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create device for teardown")
	}

	err = iota.DeleteConfigGroup(fs, resource, apiKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not create device for teardown")
	}
}
