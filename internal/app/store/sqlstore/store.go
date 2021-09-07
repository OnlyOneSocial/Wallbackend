package sqlstore

import (
	"database/sql"
	"githab.com/katelinlis/msnwallbackend/internal/app/store"
	//"githab.com/katelinlis/msnwallbackend/internal/app/model"
	_ "github.com/lib/pq" //db import
)

//Store ...
type Store struct {
	db             *sql.DB
	wallRepository *WallRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Wall() store.WallRepository {
	if s.wallRepository != nil {
		return s.wallRepository
	}

	s.wallRepository = &WallRepository{
		store: s,
	}

	return s.wallRepository
}
