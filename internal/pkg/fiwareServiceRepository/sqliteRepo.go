package fiwareservicerepository

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/google/uuid"
)

const file string = "db/sqlite/fiware.db"

type SqliteRepo struct {
	mu sync.Mutex
	db *sql.DB
}

func (sr *SqliteRepo) Add(service string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	insert, err := sr.db.Prepare(`INSERT INTO services VALUES (?,?,?,?);`)
	if err != nil {
		return err
	}
	res, err := insert.Exec(service, uuid.NewString(), time.Now(), time.Now())
	if err != nil {
		return nil
	}
	rowsEffected, err := res.RowsAffected()
	if rowsEffected == 0 {
		return fmt.Errorf("No service added")
	}
	if err != nil {
		return err
	}
	return nil
}

func (sr *SqliteRepo) Get(id string) (*FiwareServiceRow, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	get, err := sr.db.Prepare(`SELECT id, name, created_at, updated_at FROM services WHERE id = ?`)
	if err != nil {
		return nil, err
	}
	rows := get.QueryRow()
	row := FiwareServiceRow{}
	if err = rows.Scan(&row.Id, &row.Name, &row.CreatedAt, &row.UpdatedAt); err == sql.ErrNoRows {
		return nil, err
	}
	return &row, nil
}

func (sr *SqliteRepo) List() ([]FiwareServiceRow, error) {

	sr.mu.Lock()
	defer sr.mu.Unlock()
	list, err := sr.db.Prepare(`SELECT id, name, created_at, updated_at FROM services`)
	if err != nil {
		return nil, err
	}
	rows, err := list.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	services := []FiwareServiceRow{}
	for rows.Next() {
		row := FiwareServiceRow{}
		err = rows.Scan(&row.Id, &row.Name, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, err
		}
		services = append(services, row)
	}
	return services, nil

}
func (sr *SqliteRepo) Update() error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	return nil
}
func (sr *SqliteRepo) Delete(string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	return nil
}
