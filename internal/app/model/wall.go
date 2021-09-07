package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/katelinlis/goment"
)

//User ...
type Wall struct {
	ID             int    `json:"id"`
	Author         int    `json:"author"`
	AuthorUsername string `json:"author_username"`
	Text           string `json:"text"`
	Image          string //`json:"image"`
	Likes          int    //`json:"likes"`
	Timestamp      int64  //`json:"timestamp"`
	Time           string `json:"time"`
}

//Validate ...
func (w *Wall) Validate() error {
	return validation.ValidateStruct(
		w,
		validation.Field(&w.Text, validation.Required, validation.Length(1, 200)),
		validation.Field(&w.Author, validation.Required),
	)
}

func (w *Wall) Proccessing() error {

	goment.SetLocale("ru")
	time, err := goment.Unix(w.Timestamp)
	if err != nil {
		return err
	}
	w.Time = time.FromNow()

	return nil
}
