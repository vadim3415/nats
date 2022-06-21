package repository

import (
	"encoding/json"
	"fmt"
	"nats/internal/model"

	"sync"

	"github.com/jmoiron/sqlx"
)

type PqNats interface {
	PqGetId(id string) (model.ModelNats, error)
	PqNatsMsgCreate(b model.ModelNats) error
}
type Repository struct {
	PqNats
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		PqNats: NewStorage(db),
	}
}

type Storage struct {
	db  *sqlx.DB
	mtx sync.RWMutex
}

func NewStorage(db *sqlx.DB) *Storage {
	return &Storage{
		db:  db,
		mtx: sync.RWMutex{},
	}
}

///////////////////////////////////////////////////////////////////////////////
func (d *Storage) PqGetId(id string) (model.ModelNats, error) {
	d.mtx.RLock()
	defer d.mtx.RUnlock()

	var pqOutput []byte
	var output model.ModelNats

	query := fmt.Sprintf("select data from %s where order_uid=$1", usersTable)

	err := d.db.QueryRow(query, id).Scan(&pqOutput)
	if err != nil {
		return model.ModelNats{}, fmt.Errorf("task with id=%s not found %s", id, err)
	}

	if err := json.Unmarshal(pqOutput, &output); err != nil {
		return model.ModelNats{}, err
	}

	return output, nil
}

func (d *Storage) PqNatsMsgCreate(m model.ModelNats) error {
	modelMarshal, err := json.Marshal(m)
	if err != nil {
		return err
	}
	query := fmt.Sprintf("INSERT INTO %s (order_uid, data) values ($1, $2) ", usersTable)

	_, err = d.db.Exec(query, m.Order_uid, modelMarshal)
	if err != nil {
		return err
	}

	return nil
}
