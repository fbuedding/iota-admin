package fiwarerepo_test

import (
	"fmt"
	"os"
	"testing"

	//i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"

	r "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	_ "github.com/fbuedding/iota-admin/test/testing_init"
)

func TestMain(m *testing.M) {
	fmt.Println("Starting tests")
	code := m.Run()
	os.Exit(code)
}

func TestCreateSqlite(t *testing.T) {
	_, err := r.NewFiwareRepo(r.Sqlite)
	if err != nil {
		t.Error(err)
	}
}
func TestAddFiwareService(t *testing.T) {
	repo, err := r.NewFiwareRepo(r.Sqlite)
	if err != nil {
		t.Error(err)
	}
	repo.SetIdGen(func() string { return "testId" })
	err = repo.AddFiwareService("test")
	if err != nil {
		t.Error(err)
	}
}
func TestListFiwareService(t *testing.T) {
	repo, err := r.NewFiwareRepo(r.Sqlite)
	if err != nil {
		t.Error(err)
	}
	rows, err := repo.ListFiwareServices()
	if err != nil {
		t.Error(err)
	}
	if len(rows) == 0 {
		t.Log("There should be atleast one row")
		t.Fail()
	}
}

func TestUpdateFiwareService(t *testing.T) {
	repo, err := r.NewFiwareRepo(r.Sqlite)
	if err != nil {
		t.Error(err)
	}
	err = repo.UpdateFiwareService("testId", "updatedName")

	if err != nil {
		t.Error(err)
	}
}

func TestDeleteFiwareService(t *testing.T) {
	repo, err := r.NewFiwareRepo(r.Sqlite)
	if err != nil {
		t.Error(err)
	}
	err = repo.DeleteFiwareService("testId")

	if err != nil {
		t.Error(err)
	}
}
