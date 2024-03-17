package fiwareRepository

import (
	"database/sql"
	"fmt"
	"sync"
	"time"

	i "github.com/fbuedding/iota-admin/pkg/iot-agent-sdk"
	"github.com/rs/zerolog/log"
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
	rows := get.QueryRow(id)
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
		return ErrNotFound
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
		return ErrCouldNotExec
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return ErrNotFound
	}
	return nil
}

func (sr *SqliteRepo) FindFiwareServiceByName(name string) (FiwareServiceRows, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	find, err := sr.db.Prepare(`SELECT id, name, created_at, updated_at FROM services WHERE LOWER(name) LIKE CONCAT( '%',?,'%')`)
	if err != nil {
		return FiwareServiceRows{}, err
	}

	rows, err := find.Query(name)
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
	log.Debug().Any("services", services).Send()
	return services, nil
}

func (sr *SqliteRepo) AddIota(host string, port int, alias string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	insert, err := sr.db.Prepare(`INSERT INTO "iotas"("id","created_at","updated_at","port","host","alias") VALUES (?,?,?,?,?,?);`)
	if err != nil {
		return err
	}
	if alias == "" {
		alias = fmt.Sprintf("%s:%d", host, port)
	}
	res, err := insert.Exec(sr.genId(), time.Now(), time.Now(), port, host, alias)
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

func (sr *SqliteRepo) GetIota(id string) (*i.IoTA, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	get, err := sr.db.Prepare(`SELECT id, host, port, alias, created_at, updated_at FROM iotas WHERE id = ?`)
	if err != nil {
		return nil, err
	}

	rows := get.QueryRow(id)

	row := IotaRow{}
	if err = rows.Scan(&row.Id, &row.Host, &row.Port, &row.Alias, &row.CreatedAt, &row.UpdatedAt); err == sql.ErrNoRows {
		return nil, err
	}
	log.Debug().Any("IoTaRow", &row).Send()
	return row.ToIoTA(), nil
}

func (sr *SqliteRepo) ListIotas() (IotaRows, error) {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	list, err := sr.db.Prepare(`SELECT id, host, port, alias, created_at, updated_at FROM iotas`)
	if err != nil {
		return nil, err
	}
	rows, err := list.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	iotas := IotaRows{}
	for rows.Next() {
		row := IotaRow{}
		err = rows.Scan(&row.Id, &row.Host, &row.Port, &row.Alias, &row.CreatedAt, &row.UpdatedAt)
		if err != nil {
			return nil, err
		}
		iotas = append(iotas, row)
	}
	return iotas, nil
}

func (sr *SqliteRepo) DeleteIota(id string) error {
	sr.mu.Lock()
	defer sr.mu.Unlock()
	delete, err := sr.db.Prepare(`DELETE FROM iotas WHERE id = ?;`)
	if err != nil {
		return err
	}

	res, err := delete.Exec(id)
	if err != nil {
		return ErrCouldNotExec
	}
	c, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if c == 0 {
		return ErrNotFound
	}
	return nil
}
