package store

import "github.com/katelinlis/Wallbackend/internal/app/model"

//UserRepository ...
type WallRepository interface {
	Create(*model.Wall) error //Создание пользователя
	GetByAuthor(offset int, limit int, userid int) ([]model.Wall, error)
	GetByFriends(offset int, limit int, userids []int) ([]model.Wall, error)
}

type UserRepository interface {
	GetUsername(int) string //Получение имени пользователя
	GetFriends(int) []int   //Получение списка друзей пользователя
}
