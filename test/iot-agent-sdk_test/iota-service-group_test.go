package iotagentsdktest_test

import (
	"testing"
)

func TestReadConfigGroup(t *testing.T) {
	t.Log("Testing ReadConfigGroup")
	respD, err := iota.ReadConfigGroup(fs, resource, apiKey)
	if err != nil {
		t.Error(err)
	}
	if respD.Count == 0 {
		t.Fail()
	}
}

func TestListConfigGroup(t *testing.T) {
	t.Log("Testing ListConfigGroup")

	respD, err := iota.ReadConfigGroup(fs, resource, apiKey)
	if err != nil {
		t.Error(err)
	}
	if respD.Count == 0 {
		t.Fail()
	}
}

func TestUpdateConfigGroup(t *testing.T) {
	t.Log("Testing UpdateConfigGroup")
	sgtmp := sg
	sgtmp.Autoprovision = true
	err := iota.UpdateConfigGroup(fs, resource, apiKey, sgtmp)
	if err != nil {
		t.Error(err)
	}
	sgUpdated, _ := iota.ReadConfigGroup(fs, resource, apiKey)
	if sgUpdated.Services[0].Autoprovision != true {
		t.Fail()
	}
}

func TestDeleteConfigGroup(t *testing.T) {
	t.Log("Testing deleteConfigGroup")
	err := iota.DeleteConfigGroup(fs, resource, apiKey)
	if err != nil {
		t.Error(err)
	}
}

func TestUpsertConfigGroup(t *testing.T) {
	t.Log("Testing UpsertConfigGroup")

	err := iota.UpsertConfigGroup(fs, sg)
	if err != nil {
		t.Error(err)
	}
	t.Log("Testing UpsertConfigGroup again")
	err = iota.UpsertConfigGroup(fs, sg)
	if err != nil {
		t.Error(err)
	}
	_ = iota.DeleteConfigGroup(fs, resource, apiKey)
}

func TestCreatConfigGroupWSE(t *testing.T) {
	sgtemp := sg
	err := iota.CreateConfigGroupWSE(fs, &sgtemp)
	if err != nil {
		t.Error(err)
	}
	t.Log(sgtemp)
}
