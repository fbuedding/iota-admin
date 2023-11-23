package iotagentsdktest_test

import (
	"testing"
)

func TestReadServiceGroup(t *testing.T) {
	t.Log("Testing ReadServiceGroup")
	respD, err := iota.ReadServiceGroup(fs, resource, apiKey)
	if err != nil {
		t.Error(err)
	}
	if respD.Count == 0 {
		t.Fail()
	}
}
func TestListServiceGroup(t *testing.T) {
	t.Log("Testing ListServiceGroup")

	respD, err := iota.ReadServiceGroup(fs, resource, apiKey)
	if err != nil {
		t.Error(err)
	}
	if respD.Count == 0 {
		t.Fail()
	}
}

func TestUpdateServiceGroup(t *testing.T) {
	t.Log("Testing UpdateServiceGroup")
	sgtmp := sg
	sgtmp.Autoprovision = true
	err := iota.UpdateServiceGroup(fs, resource, apiKey, sgtmp)
	if err != nil {
		t.Error(err)
	}
	sgUpdated, _ := iota.ReadServiceGroup(fs, resource, apiKey)
	if sgUpdated.Services[0].Autoprovision != true {
		t.Fail()
	}
}

func TestDeleteServiceGroup(t *testing.T) {
	t.Log("Testing deleteServiceGroup")
	err := iota.DeleteServiceGroup(fs, resource, apiKey)
	if err != nil {
		t.Error(err)
	}
}

func TestUpsertServiceGroup(t *testing.T) {
	t.Log("Testing UpsertServiceGroup")

	err := iota.UpsertServiceGroup(fs, sg)
	if err != nil {
		t.Error(err)
	}
	t.Log("Testing UpsertServiceGroup again")
	err = iota.UpsertServiceGroup(fs, sg)
	if err != nil {
		t.Error(err)
	}
	iota.DeleteServiceGroup(fs, resource, apiKey)
}
func TestCreatServiceGroupWSE(t *testing.T) {
  sgtemp := sg
  err := iota.CreateServiceGroupWSE(fs, &sgtemp)
  if err != nil {
    t.Error(err)
  }
  t.Log(sgtemp)
}
