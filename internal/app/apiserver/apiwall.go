package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/katelinlis/Wallbackend/internal/app/model"
)

var ctx = context.Background()

func (s *server) ConfigureWallRouter() {

	router := s.router.PathPrefix("/api/wall").Subrouter()
	router.HandleFunc("/send", s.HandleSendWall()).Methods("POST")             // Получение всей стены
	router.HandleFunc("/get", s.HandleGetNews()).Methods("GET")                // Получение всей стены
	router.HandleFunc("/post/{postID}", s.HandleGetPost()).Methods("GET")      // Получение определенного поста
	router.HandleFunc("/get/{id}", s.HandleGetNewsByAuthor()).Methods("GET")   // Получение стены какого то пользователя
	router.HandleFunc("/like/{id}", s.HandleSetLikeOrRemove()).Methods("POST") // Поставить лайк на пост
	//router.HandleFunc("/ScanDBandCreateUUID", s.CreateUUID()).Methods("GET") //
}

// CreatePost ...
type CreatePost struct {
	Text     string `json:"text"`
	AnswerTO string `json:"answer"`
}

func (s *server) HandleSetLikeOrRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		postid, _ := vars["id"]

		userid, err := s.GetDataFromToken(w, request)
		if err != nil {
			fmt.Println(err)
			return
		}

		wall, _, err := s.store.Wall().GetPost(postid)

		if err != nil {
			return
		}

		liked, err := s.store.Wall().GetLike(wall.RandomID, int(userid.LegacyID))

		if !liked && err != nil && err.Error() == "sql: no rows in result set" {
			liked, err := s.store.Wall().SetLike(wall.RandomID, int(userid.LegacyID))
			if err != nil {
				s.error(w, request, http.StatusUnprocessableEntity, err)
				return
			}
			fmt.Println(liked)
		}

		if liked {
			liked, err := s.store.Wall().RemoveLike(wall.RandomID, int(userid.LegacyID))
			if err != nil {
				s.error(w, request, http.StatusUnprocessableEntity, err)
				return
			}
			fmt.Println(liked)
		}

	}
}

func (s *server) HandleSendWall() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		userid, err := s.GetDataFromToken(w, request)
		if err != nil {
			s.error(w, request, http.StatusUnauthorized, err)
			return
		}
		var createPost CreatePost
		err = json.NewDecoder(request.Body).Decode(&createPost)
		if err != nil {
			s.error(w, request, http.StatusBadRequest, err)
			return
		}

		var wall model.Wall
		wall.Author = int(userid.LegacyID)
		wall.Text = createPost.Text

		uuid, _ := uuid.Parse(createPost.AnswerTO)

		wall.AnswerTO = uuid
		err = s.store.Wall().Create(&wall)

		if err != nil {
			s.error(w, request, http.StatusUnprocessableEntity, err)
			return
		}

		s.redis.Del(ctx, "wallget/"+string(rune(int(userid.LegacyID))))

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		s.respond(w, request, http.StatusOK, wall)
	}
}

func (s *server) HandleGetNews() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		offset, limit := s.UrlLimitOffset(request)
		if limit > 100 {
			s.error(w, request, http.StatusBadRequest, errors.New("limit > 100"))
			return
		}

		userid, err := s.GetDataFromToken(w, request)
		if err != nil {
			s.error(w, request, http.StatusUnauthorized, err)
			return
		}

		friends, err := s.HTTPstore.User().GetFriends(int(userid.LegacyID))
		if err != nil {
			s.error(w, request, http.StatusUnprocessableEntity, err)
			return
		}
		wall, err := s.store.Wall().GetByFriends(offset, limit, friends)
		if err != nil {
			s.error(w, request, http.StatusUnprocessableEntity, err)
			return
		}
		for index, element := range wall {
			user, err := s.HTTPstore.User().GetUser(element.Author, "")
			if err != nil {
				s.error(w, request, http.StatusUnprocessableEntity, err)
				return
			}
			wall[index].AuthorUsername = user.Username
			wall[index].AuthorAvatar = user.Avatar
		}

		if err != nil {
			s.error(w, request, http.StatusUnprocessableEntity, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		s.respond(w, request, http.StatusOK, wall)
	}
}

func (s *server) HandleGetNewsByAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		authorID, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Println(err)
		}

		offset, limit := s.UrlLimitOffset(request)
		if limit > 100 {
			s.error(w, request, http.StatusBadRequest, errors.New("limit > 100"))
			return
		}

		wall, err := s.store.Wall().GetByAuthor(offset, limit, authorID)
		if err != nil {
			s.error(w, request, http.StatusUnprocessableEntity, err)
		}

		for index, element := range wall {
			user, err := s.HTTPstore.User().GetUser(element.Author, "")
			if err != nil {
				s.error(w, request, http.StatusUnprocessableEntity, err)
			}
			wall[index].AuthorUsername = user.Username
			wall[index].AuthorAvatar = user.Avatar
		}

		if err != nil {
			s.error(w, request, http.StatusUnprocessableEntity, err)
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		s.respond(w, request, http.StatusOK, wall)
	}
}

func (s *server) HandleGetPost() http.HandlerFunc {

	type PostData struct {
		Post    model.Wall
		Answers []model.Wall
	}

	return func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		postID := vars["postID"]

		post, answers, err := s.store.Wall().GetPost(postID)
		if err != nil {
			s.error(w, request, http.StatusUnprocessableEntity, err)
			return
		}

		user, err := s.HTTPstore.User().GetUser(post.Author, "")
		if err != nil {
			s.error(w, request, http.StatusUnprocessableEntity, err)
			return
		}
		post.AuthorUsername = user.Username
		post.AuthorAvatar = user.Avatar

		for index, element := range answers {
			user, err := s.HTTPstore.User().GetUser(element.Author, "")
			if err != nil {
				s.error(w, request, http.StatusUnprocessableEntity, err)
				return
			}
			answers[index].AuthorUsername = user.Username
			answers[index].AuthorAvatar = user.Avatar
		}

		postdata := PostData{
			Post:    post,
			Answers: answers,
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		s.respond(w, request, http.StatusOK, postdata)
	}
}

func (s *server) CreateUUID() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		s.store.Wall().ScanAndCreateUUID()
		s.respond(w, request, http.StatusOK, "ok")
	}
}
