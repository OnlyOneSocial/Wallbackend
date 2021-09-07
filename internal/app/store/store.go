package store

/*
Store репозитории данных
*/
type Store interface {
	Wall() WallRepository // интерфейс для стены
}

type HTTPStore interface {
	User() UserRepository //интерфейс для пользователей
}
