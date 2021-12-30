package httpstore

import (
	"time"

	"github.com/katelinlis/Wallbackend/internal/app/store"
	"github.com/patrickmn/go-cache"
)

//Store ...
type Store struct {
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
	return &Store{
		CacheUser:    cache.New(5*time.Minute, 10*time.Minute),
		CacheFriends: cache.New(5*time.Minute, 10*time.Minute),
		userCache:    make(map[int]string),
		friendsCache: make(map[int][]int),
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
