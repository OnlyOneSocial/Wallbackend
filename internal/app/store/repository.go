package store

import (
	"github.com/google/uuid"
	"github.com/katelinlis/Wallbackend/internal/app/model"
)

//WallRepository ...
type WallRepository interface {
	Create(*model.Wall) error                                                //Создание пользователя
	GetByAuthor(offset int, limit int, userid int) ([]model.Wall, error)     // Получить новости за определенного пользователя
	GetByFriends(offset int, limit int, userids []int) ([]model.Wall, error) // Получить новости друзей и людей на которых подписан пользователь
	GetPost(PostID string) (model.Wall, []model.Wall, error)                 // Получение определенного поста
	ScanAndCreateUUID() error                                                // Сканирование и создание UUID если пусто
	GetAnswersCount(PostID string) (int, error)
	GetAnswers(PostID string) ([]model.Wall, error)
	SetLike(PostID uuid.UUID, wholiked int) (bool, error)
	RemoveLike(PostID uuid.UUID, wholiked int) (bool, error)
	GetLike(PostID uuid.UUID, wholiked int) (bool, error)
}

//UserRepository ...
type UserRepository interface {
	GetUser(AuthorID int) (model.UserObj, error) //Получение данных о пользователе
	GetFriends(int) ([]int, error)               //Получение списка друзей пользователя
}
