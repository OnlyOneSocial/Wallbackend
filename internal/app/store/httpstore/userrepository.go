package httpstore

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/katelinlis/Wallbackend/internal/app/model"
	"github.com/patrickmn/go-cache"
)

//UserRepository ...
type UserRepository struct {
	store *Store
}

//GetUser ...
func (r *UserRepository) GetUser(AuthorID int) model.UserObj {
	if val, ok := r.store.CacheUser.Get(strconv.Itoa(AuthorID)); ok {
		return val.(model.UserObj)
	}
	userID := strconv.Itoa(AuthorID)

	client := http.Client{}
	resp, err := client.Get(`http://localhost:3046/api/user/get/` + userID)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	usrObj := model.UserObj{
		Username: result["user"]["username"],
		Avatar:   result["user"]["avatar"],
	}
	r.store.CacheUser.Set(userID, usrObj, cache.DefaultExpiration)
	return usrObj
}

//GetFriends ...
func (r *UserRepository) GetFriends(AuthorID int) []int {
	if val, ok := r.store.CacheFriends.Get(strconv.Itoa(AuthorID)); ok {
		return val.([]int)
	}
	userID := strconv.Itoa(AuthorID)

	client := http.Client{}
	resp, err := client.Get(`http://localhost:3046/api/friends/array_friends/` + userID)
	if err != nil {
		log.Fatalln(err)
	}

	var result []int
	json.NewDecoder(resp.Body).Decode(&result)

	r.store.CacheFriends.Set(userID, result, cache.DefaultExpiration)
	return result
}
