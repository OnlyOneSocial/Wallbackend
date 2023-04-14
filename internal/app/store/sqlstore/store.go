package sqlstore

import (
	"database/sql"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/katelinlis/Wallbackend/internal/app/store"

	//"githab.com/katelinlis/msnwallbackend/internal/app/model"
	_ "github.com/lib/pq" //db import
)

// Store ...
type Store struct {
	cache          *cache.Cache
	db             *sql.DB
	wallRepository *WallRepository
}

// New ...
func New(db *sql.DB) *Store {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	mycache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(10000, time.Minute),
	})

	return &Store{
		db:    db,
		cache: mycache,
	}
}

// Wall ...
func (s *Store) Wall() store.WallRepository {
	if s.wallRepository != nil {
		return s.wallRepository
	}

	s.wallRepository = &WallRepository{
		store: s,
	}

	return s.wallRepository
}
