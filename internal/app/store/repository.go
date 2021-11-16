package store

import "github.com/katelinlis/Wallbackend/internal/app/model"

//WallRepository ...
type WallRepository interface {
	Create(*model.Wall) error                                                //Создание пользователя
	GetByAuthor(offset int, limit int, userid int) ([]model.Wall, error)     // Получить новости за определенного пользователя
	GetByFriends(offset int, limit int, userids []int) ([]model.Wall, error) // Получить новости друзей и людей на которых подписан пользователь
	GetPost(PostID string) (model.Wall, []model.Wall, error)                 // Получение определенного поста
	ScanAndCreateUUID() error                                                // Сканирование и создание UUID если пусто
}

//UserRepository ...
type UserRepository interface {
	GetUsername(int) string //Получение имени пользователя
	GetFriends(int) []int   //Получение списка друзей пользователя
}
