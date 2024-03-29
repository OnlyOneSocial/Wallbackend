package sqlstore

import (
	//"database/sql"

	"fmt"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/google/uuid"
	"github.com/katelinlis/Wallbackend/internal/app/model"
	"github.com/lib/pq"
)

// WallRepository ...
type WallRepository struct {
	store *Store
}

// Create ...
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

// Update ...
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

// SetLike ...
func (r *WallRepository) SetLike(PostID uuid.UUID, wholiked int) (bool, error) {
	err := r.store.db.QueryRow("UPDATE wall SET likes = array_append(likes,$1) WHERE random_id=$2 RETURNING random_id", wholiked, PostID.String())

	var id []uint8
	if err := err.Scan(&id); err != nil {
		return false, err
	}

	return true, nil
}

// RemoveLike ...
func (r *WallRepository) RemoveLike(PostID uuid.UUID, wholiked int) (bool, error) {
	err := r.store.db.QueryRow("UPDATE wall SET likes = array_remove(likes,$1) WHERE random_id=$2 RETURNING random_id", wholiked, PostID.String())

	var id []uint8
	if err := err.Scan(&id); err != nil {
		return false, err
	}

	return true, nil
}

// GetLike ...
func (r *WallRepository) GetLike(PostID uuid.UUID, wholiked int) (bool, error) {
	err := r.store.db.QueryRow("Select random_id from wall where random_id = $2 and $1 = ANY(likes)", wholiked, PostID.String())

	var id []uint8
	if err := err.Scan(&id); err != nil {
		return false, err
	}

	return true, nil
}

// GetLikes ...
func (r *WallRepository) GetLikes(PostID uuid.UUID) (bool, error) {
	err := r.store.db.QueryRow("Select random_id,likes::bigint[] from wall where random_id = $1", PostID.String())

	var id []uint8
	if err := err.Scan(&id); err != nil {
		return false, err
	}

	return true, nil
}

// GetByAuthor ...
func (r *WallRepository) GetByAuthor(offset int, limit int, userid int) ([]model.Wall, error) {
	wall := []model.Wall{}

	err := r.store.cache.Once(&cache.Item{
		Key:   "newsByAuthor" + fmt.Sprint(userid) + fmt.Sprint(limit) + fmt.Sprint(offset),
		Value: &wall, // destination
		TTL:   time.Second * 5,
		Do: func(*cache.Item) (interface{}, error) {

			localWall := []model.Wall{}

			var rows, err = r.store.db.Query("Select author,text,timestamp,random_id,answer_to,likes::bigint[] from wall where author = $1 ORDER BY timestamp DESC limit $2 OFFSET $3", userid, limit, offset)
			for rows.Next() {
				post := model.Wall{}
				err := rows.Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID, &post.AnswerTO, pq.Array(&post.Likes))
				if err != nil {
					return localWall, err
				}
				post.Proccessing()
				AnswerCount, err := r.GetAnswersCount(post.RandomID.String())
				if err != nil {
					return localWall, err
				}
				post.AnswerCount = AnswerCount
				localWall = append(localWall, post)
			}
			fmt.Println(err)
			return localWall, err
		},
	})
	fmt.Println(err)
	return wall, err
}

// GetByFriends ...
func (r *WallRepository) GetByFriends(offset int, limit int, userids []int) ([]model.Wall, error) {
	wall := []model.Wall{}

	var rows, err2 = r.store.db.Query("select author,text,timestamp,random_id,answer_to,likes::bigint[] from wall where author = ANY($1::int[]) ORDER BY timestamp DESC limit $2 OFFSET $3", pq.Array(userids), limit, offset)

	if err2 != nil {
		return wall, err2
	}

	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID, &post.AnswerTO, pq.Array(&post.Likes))
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

// GetAnswers ...
func (r *WallRepository) GetAnswers(PostID string) ([]model.Wall, error) {
	answer := []model.Wall{}
	rows, err := r.store.db.Query(
		"select author,text,timestamp,random_id,answer_to,likes::bigint[] from wall where answer_to = $1",
		PostID,
	)
	if err != nil {
		return answer, err
	}

	for rows.Next() {
		post := model.Wall{}
		err := rows.Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID, &post.AnswerTO, pq.Array(&post.Likes))
		if err != nil {
			return answer, err
		}
		AnswerCount, err := r.GetAnswersCount(post.RandomID.String())
		if err != nil {
			return answer, err
		}
		post.AnswerCount = AnswerCount
		post.Proccessing()

		answer = append(answer, post)

	}
	return answer, err
}

// GetAnswersCount ...
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

// GetPost ...
func (r *WallRepository) GetPost(PostID string) (model.Wall, []model.Wall, error) {
	post := model.Wall{}
	answer := []model.Wall{}
	var likes pq.Int64Array
	var err = r.store.db.QueryRow(
		"select author,text,timestamp,random_id,answer_to,likes::bigint[] from wall where random_id = $1",
		PostID,
	).Scan(&post.Author, &post.Text, &post.Timestamp, &post.RandomID, &post.AnswerTO, &likes)
	if err != nil {
		return post, answer, err
	}
	post.Likes = []int64(likes)

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
