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
	CacheUser      *cache.Cache
	CacheFriends   *cache.Cache
	userCache      map[int]string
	friendsCache   map[int][]int
	userRepository *UserRepository
}

/*
Todo поставить ограничение на карту

 и удалять некоторую информацию в случае достяжении лимита

 for k := range userCache {
    delete(userCache, k)
}

а лучше подключить это https://github.com/patrickmn/go-cache
*/

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

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
	}

	return s.userRepository
}
