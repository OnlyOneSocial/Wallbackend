package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/google/uuid"
	"github.com/katelinlis/goment"
)

//Wall ...
type Wall struct {
	ID             int       `json:"id"`
	Author         int       `json:"author"`
	AuthorUsername string    `json:"author_username"`
	Text           string    `json:"text"`
	Image          string    //`json:"image"`
	Likes          int       //`json:"likes"`
	Timestamp      int64     //`json:"timestamp"`
	Time           string    `json:"time"`
	RandomID       uuid.UUID `json:"random_id"`
	AnswerTO       uuid.UUID `json:"answerto"`
}

//Validate ...
func (w *Wall) Validate() error {
	return validation.ValidateStruct(
		w,
		validation.Field(&w.Text, validation.Required, validation.Length(1, 400)),
		validation.Field(&w.Author, validation.Required),
	)
}

//Proccessing ...
func (w *Wall) Proccessing() error {

	goment.SetLocale("ru")
	time, err := goment.Unix(w.Timestamp)
	if err != nil {
		return err
	}
	w.Time = time.FromNow()

	return nil
}

//GenerateUUID ...
func (w *Wall) GenerateUUID() error {
	uuid, err := uuid.NewRandom()
	if err != nil {
		return err
	}
	w.RandomID = uuid
	return nil
}
