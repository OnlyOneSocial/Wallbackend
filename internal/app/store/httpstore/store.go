package httpstore

import (
	"github.com/katelinlis/Wallbackend/internal/app/store"
)

//Store ...
type Store struct {
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
