package fiwareservicerepository

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"
	"time"
)

const file string = "db/sqlite/fiware.db"

type SqliteRepo struct {
	mu    sync.Mutex
	db    *sql.DB
	genId func() string
}

func (sr *SqliteRepo) SetIdGen(f func() string) {
	sr.genId = f
}

func (sr *SqliteRepo) AddFiwareService(service string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	insert, err := sr.db.Prepare(`INSERT INTO services VALUES (?,?,?,?);`)
	if err != nil {
		return err
	}
	res, err := insert.Exec(service, sr.genId(), time.Now(), time.Now())
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

func (sr *SqliteRepo) GetFiwareService(id string) (*FiwareServiceRow, error) {
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

func (sr *SqliteRepo) ListFiwareServices() (FiwareServiceRows, error) {

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
	services := FiwareServiceRows{}
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
func (sr *SqliteRepo) UpdateFiwareService(id string, name string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	update, err := sr.db.Prepare(`UPDATE services SET name = ?, updated_at = ? WHERE id = ?;`)

	if err != nil {
		return err
	}

	res, err := update.Exec(name, time.Now(), id)

	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New("No row affected")
	}
	return nil
}
func (sr *SqliteRepo) DeleteFiwareService(id string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	delete, err := sr.db.Prepare(`DELETE FROM services WHERE id = ?;`)

	if err != nil {
		return err
	}

	res, err := delete.Exec(id)

	if err != nil {
		return err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return errors.New("No row affected")
	}
	return nil
}
