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

func TestInitIntegrations(t *testing.T) {
    t.Log("TestMigrations")
	is := r.GetMigrations()
	if len(is) == 0 {
    t.Log("No migrations found")
		t.Fail()
	}
}

func TestCreateSqlite(t *testing.T) {
	_, err := r.NewFiwareServiceRepo(r.Sqlite)
	if err != nil {
		t.Error(err)
	}
}
func TestAdd(t *testing.T) {
	repo, err := r.NewFiwareServiceRepo(r.Sqlite)
	if err != nil {
		t.Error(err)
	}
	err = repo.AddFiwareService("test")
	if err != nil {
		t.Error(err)
	}
}
func TestList(t *testing.T) {
	repo, err := r.NewFiwareServiceRepo(r.Sqlite)
	if err != nil {
		t.Error(err)
	}
	_, err = repo.ListFiwareServices()
	if err != nil {
		t.Error(err)
	}
}
