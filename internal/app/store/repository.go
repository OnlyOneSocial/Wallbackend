package store

import "githab.com/katelinlis/msnwallbackend/internal/app/model"

//UserRepository ...
type WallRepository interface {
	Create(*model.Wall) error //Создание пользователя
	GetByAuthor(offset int, limit int, userid int) ([]model.Wall, error)
	GetByFriends(offset int, limit int, userids string) ([]model.Wall, error)
}

type UserRepository interface {
	GetUsername(int) string //Получение имени пользователя
}
