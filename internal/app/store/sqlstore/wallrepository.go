package sqlstore

import (
	//"database/sql"

	"time"

	"github.com/katelinlis/Wallbackend/internal/app/model"
	"github.com/lib/pq"
)

//WallRepository ...
type WallRepository struct {
	store *Store
}

//Create ...
func (r *WallRepository) Create(p *model.Wall) error {

	if err := p.Validate(); err != nil {
		return err
	}
	p.GenerateUUID()

	var err2 = r.store.db.QueryRow(
		"INSERT INTO wall (author,text,timestamp,random_id) VALUES ($1,$2,$3,$4) RETURNING id",
		p.Author,
		p.Text,
		time.Now().Unix(),
		p.RandomID.String(),
	).Scan(&p.ID)

	return err2
}

//Update ...
func (r *WallRepository) Update(p *model.Wall) error {

	if err := p.Validate(); err != nil {
		return err
	}

	var err2 = r.store.db.QueryRow(
		"UPDATE wall set random_id = $1 where id=$2 RETURNING id",
		p.RandomID.String(),
		p.ID,
	).Scan(&p.ID)

	return err2
}

//GetByAuthor ...
func (r *WallRepository) GetByAuthor(offset int, limit int, userid int) ([]model.Wall, error) {
	wall := []model.Wall{}

	var rows, err2 = r.store.db.Query("Select author,text,timestamp,random_id from wall where author = $1 ORDER BY id DESC limit $2 OFFSET $3", userid, limit, offset)
	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID)
		if err != nil {
			return wall, err
		}
		post.Proccessing()
		wall = append(wall, post)
	}

	return wall, err2
}

// GetByFriends ...
func (r *WallRepository) GetByFriends(offset int, limit int, userids []int) ([]model.Wall, error) {
	wall := []model.Wall{}

	var rows, err2 = r.store.db.Query("select author,text,timestamp,random_id from wall where author = ANY($1::int[]) ORDER BY id DESC limit $2 OFFSET $3", pq.Array(userids), limit, offset)

	if err2 != nil {
		return wall, err2
	}

	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID)
		if err != nil {
			return wall, err
		}
		post.Proccessing()

		wall = append(wall, post)
	}

	return wall, err2
}

//GetPost ...
func (r *WallRepository) GetPost(AuthorID int, PostID string) (model.Wall, error) {
	post := model.Wall{}

	var err = r.store.db.QueryRow(
		"select author,text,timestamp,random_id from wall where random_id = $1 and author = $2",
		PostID,
		AuthorID,
	).Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID)

	post.Proccessing()

	//post.RandomID = uuid
	return post, err
}

// ScanAndCreateUUID ...
func (r *WallRepository) ScanAndCreateUUID() error {
	var rows, err = r.store.db.Query("select id,author,text,timestamp from wall where random_id IS NULL ORDER BY id limit 100")

	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.ID, &post.Author, &post.Text, &post.Timestamp)
		if err != nil {
			return err
		}
		post.GenerateUUID()
		r.Update(&post)
	}
	return err
}
