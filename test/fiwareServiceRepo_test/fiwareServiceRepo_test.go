package fiwareservicerepotest_test

import (
	"os"
	"testing"

	//i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	r "github.com/fbuedding/iota-admin/internal/pkg/fiwareServiceRepository"
	_ "github.com/fbuedding/iota-admin/test/testing_init"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	d, _ := os.Getwd()
	log.Debug().Str("workingDirect", d).Msg("Init Fiware Service Repo tests")
}

func TestInitIntegrations(t *testing.T) {
	is := r.GetIntegrations()
	if len(is) == 0 {
		t.Fail()
	}
}

func TestCreateSqlite(t *testing.T) {
	_, err := r.NewFiwareServiceRepo(r.Sqlite)
	if err != nil {
		t.Fail()
	}
}
func TestAdd(t *testing.T) {
	repo, err := r.NewFiwareServiceRepo(r.Sqlite)
	if err != nil {
		t.Fail()
	}
	err = repo.Add("test")
	if err != nil {
		t.Fail()
	}
}
func TestList(t *testing.T) {
	repo, err := r.NewFiwareServiceRepo(r.Sqlite)
	if err != nil {
		t.Fail()
	}
	_, err = repo.List()
	if err != nil {
		log.Error().Err(err).Send()
		t.Fail()
	}
}
