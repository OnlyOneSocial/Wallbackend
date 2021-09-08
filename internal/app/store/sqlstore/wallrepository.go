package sqlstore

import (
	//"database/sql"

	"time"

	"github.com/katelinlis/Wallbackend/internal/app/model"
)

type WallRepository struct {
	store *Store
}

func (r *WallRepository) Create(p *model.Wall) error {

	if err := p.Validate(); err != nil {
		return err
	}

	var err2 = r.store.db.QueryRow(
		"INSERT INTO wall (author,text,timestamp) VALUES ($1,$2,$3) RETURNING id",
		p.Author,
		p.Text,
		time.Now().Unix(),
	).Scan(&p.ID)

	return err2

}

func (r *WallRepository) GetByAuthor(offset int, limit int, userid int) ([]model.Wall, error) {
	wall := []model.Wall{}

	var rows, err2 = r.store.db.Query("Select author,text,timestamp from wall where author = $1 ORDER BY id DESC limit $2 OFFSET $3", userid, limit, offset)
	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.Author, &post.Text, &post.Timestamp)
		if err != nil {
			return wall, err
		}
		post.Proccessing()
		wall = append(wall, post)
	}

	return wall, err2
}

func (r *WallRepository) GetByFriends(offset int, limit int, userids string) ([]model.Wall, error) {
	wall := []model.Wall{}

	var rows, err2 = r.store.db.Query("select author,text,timestamp from wall where author in ($1) ORDER BY id DESC limit $2 OFFSET $3", userids, limit, offset)
	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.Author, &post.Text, &post.Timestamp)
		if err != nil {
			return wall, err
		}
		post.Proccessing()

		wall = append(wall, post)
	}

	return wall, err2
}
