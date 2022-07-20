package httpstore

import (
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/katelinlis/Wallbackend/internal/app/store"
)

//Store ...
type Store struct {
	cache          *cache.Cache
	userRepository *UserRepository
}

//New ...
func New() *Store {

	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	mycache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(10000, time.Minute),
	})

	return &Store{
		cache: mycache,
	}
}

//User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
