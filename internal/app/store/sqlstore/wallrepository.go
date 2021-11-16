package sqlstore

import (
	//"database/sql"

	"fmt"
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
		"INSERT INTO wall (author,text,timestamp,random_id,answer_to) VALUES ($1,$2,$3,$4,$5) RETURNING id",
		p.Author,
		p.Text,
		time.Now().Unix(),
		p.RandomID.String(),
		p.AnswerTO.String(),
	).Scan(&p.ID)

	return err2
}

//Update ...
func (r *WallRepository) Update(p *model.Wall) error {

	if err := p.Validate(); err != nil {
		return err
	}

	fmt.Println(p.RandomID.String())

	var err2 = r.store.db.QueryRow(
		"UPDATE wall set answer_to = $1 where random_id=$2 RETURNING id",
		p.AnswerTO.String(),
		p.RandomID.String(),
	).Scan(&p.ID)

	return err2
}

//GetByAuthor ...
func (r *WallRepository) GetByAuthor(offset int, limit int, userid int) ([]model.Wall, error) {
	wall := []model.Wall{}

	var rows, err2 = r.store.db.Query("Select author,text,timestamp,random_id,answer_to from wall where author = $1 ORDER BY id DESC limit $2 OFFSET $3", userid, limit, offset)
	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID, &post.AnswerTO)
		if err != nil {
			return wall, err
		}
		post.Proccessing()
		AnswerCount, err := r.GetAnswersCount(post.RandomID.String())
		if err != nil {
			return wall, err
		}
		post.AnswerCount = AnswerCount
		wall = append(wall, post)
	}

	return wall, err2
}

// GetByFriends ...
func (r *WallRepository) GetByFriends(offset int, limit int, userids []int) ([]model.Wall, error) {
	wall := []model.Wall{}

	var rows, err2 = r.store.db.Query("select author,text,timestamp,random_id,answer_to from wall where author = ANY($1::int[]) ORDER BY id DESC limit $2 OFFSET $3", pq.Array(userids), limit, offset)

	if err2 != nil {
		return wall, err2
	}

	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID, &post.AnswerTO)
		if err != nil {
			return wall, err
		}
		AnswerCount, err := r.GetAnswersCount(post.RandomID.String())
		if err != nil {
			return wall, err
		}
		post.AnswerCount = AnswerCount
		post.Proccessing()

		wall = append(wall, post)
	}

	return wall, err2
}

//GetAnswers ...
func (r *WallRepository) GetAnswers(PostID string) ([]model.Wall, error) {
	answer := []model.Wall{}
	rows, err := r.store.db.Query(
		"select author,text,timestamp,random_id,answer_to from wall where answer_to = $1",
		PostID,
	)
	if err != nil {
		return answer, err
	}

	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID, &post.AnswerTO)
		if err != nil {
			return answer, err
		}
		AnswerCount, err := r.GetAnswersCount(PostID)
		if err != nil {
			return answer, err
		}
		post.AnswerCount = AnswerCount
		post.Proccessing()

		answer = append(answer, post)

	}
	return answer, err
}

//GetAnswersCount ...
func (r *WallRepository) GetAnswersCount(PostID string) (int, error) {
	var count int
	err := r.store.db.QueryRow(
		"select COUNT(answer_to) from wall where answer_to = $1",
		PostID,
	).Scan(&count)
	if err != nil {
		return count, err
	}
	return count, err
}

//GetPost ...
func (r *WallRepository) GetPost(PostID string) (model.Wall, []model.Wall, error) {
	post := model.Wall{}
	answer := []model.Wall{}

	var err = r.store.db.QueryRow(
		"select author,text,timestamp,random_id,answer_to from wall where random_id = $1",
		PostID,
	).Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID, &post.AnswerTO)
	if err != nil {
		return post, answer, err
	}

	answer, err = r.GetAnswers(PostID)
	if err != nil {
		return post, answer, err
	}
	AnswerCount, err := r.GetAnswersCount(post.RandomID.String())
	if err != nil {
		return post, answer, err
	}
	post.AnswerCount = AnswerCount
	post.Proccessing()

	//post.RandomID = uuid
	return post, answer, err
}

// ScanAndCreateUUID ...
func (r *WallRepository) ScanAndCreateUUID() error {
	var rows, err = r.store.db.Query("select random_id,author,text,timestamp from wall where answer_to IS NULL ORDER BY id limit 100")

	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.RandomID, &post.Author, &post.Text, &post.Timestamp)
		if err != nil {
			return err
		}
		//post.GenerateUUID()
		r.Update(&post)
	}
	return err
}
