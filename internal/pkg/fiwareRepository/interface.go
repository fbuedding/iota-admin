package fiwareRepository

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"time"

	i "github.com/fbuedding/fiware-iot-agent-sdk"
	"github.com/google/uuid"

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

type IotaRows []IotaRow

type IotaRow struct {
	Id        string
	Alias     string
	Host      string
	Port      int
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
		ServicePath: "/*",
	}
}

func (rs IotaRows) ToIoTAs() []*i.IoTA {
	fss := make([]*i.IoTA, len(rs))
	for i, v := range rs {
		fss[i] = v.ToIoTA()
	}
	return fss
}

func (r IotaRow) ToIoTA() *i.IoTA {
	return &i.IoTA{
		Host: r.Host,
		Port: r.Port,
	}
}

type FiwareRepo interface {
	AddFiwareService(string) error
	GetFiwareService(string) (*FiwareServiceRow, error)
	ListFiwareServices() (FiwareServiceRows, error)
	UpdateFiwareService(string, string) error
	DeleteFiwareService(string) error
	SetIdGen(func() string)
	FindFiwareServiceByName(string) (FiwareServiceRows, error)
	AddIota(string, int, string) error
	GetIota(string) (*i.IoTA, error)
	ListIotas() (IotaRows, error)
	DeleteIota(string) error
}

const (
	Sqlite = iota
)

type RepoType int

func (r RepoType) String() string {
	return [...]string{"Sqlite"}[r]
}

var (
	ErrNotFound      = errors.New("not found")
	ErrCouldNotQuery = errors.New("couldn't query")
	ErrCouldNotExec  = errors.New("couldn't Exec")
)

func init() {
}

func NewFiwareRepo(i RepoType) (FiwareRepo, error) {
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

	repo := &SqliteRepo{db: db, genId: uuid.NewString}
	err = repo.initRepo()
	if err != nil {
		return nil, err
	}

	return repo, nil

}
