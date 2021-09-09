package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/katelinlis/Wallbackend/internal/app/model"
)

func (s *server) ConfigureWallRouter() {

	router := s.router.PathPrefix("/api/wall").Subrouter()
	router.HandleFunc("/send", s.HandleSendWall()).Methods("POST")           // Получение всей стены
	router.HandleFunc("/get", s.HandleGetNews()).Methods("GET")              // Получение всей стены
	router.HandleFunc("/get/{id}", s.HandleGetNewsByAuthor()).Methods("GET") // Получение стены какого то пользователя
}

type CreatePost struct {
	Text string `json:"text"`
}

func (s *server) HandleSendWall() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		userid, err := s.GetDataFromToken(w, request)
		if err != nil {
			fmt.Println(err)
		}
		var createPost CreatePost
		json.NewDecoder(request.Body).Decode(&createPost)

		var wall model.Wall
		wall.Author = int(userid)
		wall.Text = createPost.Text
		err = s.store.Wall().Create(&wall)

		if err != nil {
			s.respond(w, request, http.StatusUnprocessableEntity, err)
			return
		}

		s.redis.Del("wallget/" + string(rune(int(userid))))
		s.respond(w, request, http.StatusOK, wall)
	}
}

func (s *server) HandleGetNews() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		offset, limit := s.UrlLimitOffset(request)

		userid, err := s.GetDataFromToken(w, request)
		if err != nil {
			fmt.Println(err)
		}

		friends := s.HTTPstore.User().GetFriends(int(userid))
		wall, err := s.store.Wall().GetByFriends(offset, limit, friends)
		if err != nil {
			fmt.Println(err)
		}
		for index, element := range wall {
			username := s.HTTPstore.User().GetUsername(element.Author)
			wall[index].AuthorUsername = username
		}
		s.respond(w, request, http.StatusOK, wall)
	}
}

func (s *server) HandleGetNewsByAuthor() http.HandlerFunc {
	return func(w http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			fmt.Println(err)
		}

		offset, limit := s.UrlLimitOffset(request)
		wall, err := s.store.Wall().GetByAuthor(offset, limit, id)
		if err != nil {
			fmt.Println(err)
		}

		for index, element := range wall {
			username := s.HTTPstore.User().GetUsername(element.Author)
			wall[index].AuthorUsername = username
		}

		s.respond(w, request, http.StatusOK, wall)
	}
}
