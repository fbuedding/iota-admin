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

type FiwareServiceRows []FiwareServiceRow

type FiwareServiceRow struct {
	Id        string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (rs FiwareServiceRows) ToFiwareServices() []*i.FiwareService {
	fss := make([]*i.FiwareService, len(rs))
	for i, v := range rs {
		fss[i] = v.ToFiwareService()
	}
	return fss
}

func (r FiwareServiceRow) ToFiwareService() *i.FiwareService {
	return &i.FiwareService{
		Service:     r.Name,
		ServicePath: "/#",
	}
}

type FiwareRepo interface {
	AddFiwareService(string) error
	GetFiwareService(string) (*FiwareServiceRow, error)
	ListFiwareServices() (FiwareServiceRows, error)
	UpdateFiwareService() error
	DeleteFiwareService(string) error
}

const (
	Sqlite = iota
)

type RepoType int

func (r RepoType) String() string {
	return [...]string{"Sqlite"}[r]
}

var (
	//go:embed migrations/*.sql
	migrations_files embed.FS
	migrations       []string
)

func init() {
	migrations = make([]string, 0)
	files, err := migrations_files.ReadDir("migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("Could not read files")
	}
	for _, f := range files {
		log.Debug().Msg(f.Name())

		data, err := migrations_files.ReadFile("migrations/" + f.Name())
		if err != nil {
			fmt.Println(err)
			log.Fatal().Str("file", f.Name()).Msg("Could not open file")
		}
		migrations = append(migrations, string(data))
	}
}

func GetMigrations() []string {
	return migrations
}

func NewFiwareServiceRepo(i RepoType) (FiwareRepo, error) {
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
	for _, v := range migrations {
		_, err := db.Exec(v)
		if err != nil {
			return nil, err
		}
	}
	return &SqliteRepo{db: db}, nil
}
