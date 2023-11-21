package fiwareservicerepository

import (
	"database/sql"
	"embed"
	_ "embed"
	"errors"
	"fmt"
	"time"

	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"

	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/zerolog/log"
)
type FiwareServiceRow struct {
  Id string
  Name string
  CreatedAt time.Time
  UpdatedAt time.Time
}

func (r FiwareServiceRow)ToFiwareService()*i.FiwareService{
  return &i.FiwareService{
  	Service:    r.Name,
  	ServicePath: "/#",
  }
}

type FiwareServiceRepo interface {
	Add(string) error
	Get(string) (*FiwareServiceRow, error)
	List() ([]FiwareServiceRow, error)
	Update() error
	Delete(string) error
}


const (
	Sqlite = iota
)

type RepoType int

func (r RepoType) String() string {
	return [...]string{"Sqlite"}[r]
}

var (
	//go:embed integrations/*.sql
	integrations_files embed.FS
	integrations       []string
)

func init() {
	integrations = make([]string, 0)
	files, err := integrations_files.ReadDir("integrations")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not read files")
	}
	for _, f := range files {
		log.Debug().Msg(f.Name())

		data, err := integrations_files.ReadFile("integrations/" + f.Name())
		if err != nil {
			fmt.Println(err)
			log.Fatal().Str("file", f.Name()).Msg("Could not open file")
		}
		integrations = append(integrations, string(data))
	}
}

func GetIntegrations() []string {
	return integrations
}

func NewFiwareServiceRepo(i RepoType) (FiwareServiceRepo, error) {
	switch i {
	case Sqlite:
		repo, err := newSqliteRepo()
		if err != nil {
			log.Fatal().Err(err).Msg("Could not connecto to db")
			return nil, err
		}
		return repo, nil
	default:
		return nil, errors.New("invalid repo")
	}
}

func newSqliteRepo() (*SqliteRepo, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return nil, fmt.Errorf("Could not open Sqlite db: %w", err)
	}
	for _, v := range integrations {
    _, err := db.Exec(v)
    if err != nil {
      return nil, err
    }
	}
	return &SqliteRepo{db: db}, nil
}
