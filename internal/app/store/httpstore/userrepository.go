package httpstore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) GetUsername(AuthorID int) string {
	if val, ok := r.store.userCache[AuthorID]; ok {
		return val
	}
	userId := strconv.Itoa(AuthorID)

	client := http.Client{}
	resp, err := client.Get(`http://localhost:3044/api/user/get/` + userId)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]map[string]string
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println(result["user"]["username"])
	r.store.userCache[AuthorID] = result["user"]["username"]
	return result["user"]["username"]
}

func (r *UserRepository) GetFriends(AuthorID int) []int {
	if val, ok := r.store.friendsCache[AuthorID]; ok {
		return val
	}
	userId := strconv.Itoa(AuthorID)

	client := http.Client{}
	resp, err := client.Get(`http://localhost:3044/api/user/array_friends/` + userId)
	if err != nil {
		log.Fatalln(err)
	}

	var result []int
	json.NewDecoder(resp.Body).Decode(&result)

	fmt.Println(result)
	r.store.friendsCache[AuthorID] = result
	return result
}