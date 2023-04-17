package httpstore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/katelinlis/Wallbackend/internal/app/model"
)

//UserRepository ...
type UserRepository struct {
	store *Store
}

//GetUser ...
func (r *UserRepository) GetUser(AuthorID int, useruuid string) (usrObj model.UserObj, err error) {
	/*if val, ok := r.store.CacheUser.Get(strconv.Itoa(AuthorID)); ok {
		return val.(model.UserObj)
	}*/
	userID := strconv.Itoa(AuthorID)

	err = r.store.cache.Once(&cache.Item{
		Key:   "user" + useruuid + fmt.Sprint(userID),
		Value: &usrObj, // destination
		TTL:   time.Hour,
		Do: func(*cache.Item) (interface{}, error) {

			client := http.Client{}
			var resp *http.Response
			if useruuid == "" {
				resp, err = client.Get(`https://only-one.su/api/user/get/` + userID)
			} else {
				resp, err = client.Get(`https://only-one.su/api/user/get_by_uuid/` + useruuid)
			}
			if err != nil {
				log.Fatalln(err)
			}

			var result map[string]map[string]string
			json.NewDecoder(resp.Body).Decode(&result)
			userIDInt, err := strconv.Atoi(result["user"]["id"])
			usrObj := model.UserObj{
				ID:       userIDInt,
				Username: result["user"]["username"],
				Avatar:   result["user"]["avatar"],
			}

			return usrObj, err
		},
	})
	return usrObj, err
}

//GetFriends ...
func (r *UserRepository) GetFriends(AuthorID int) (result []int, err error) {

	err = r.store.cache.Once(&cache.Item{
		Key:   "array_friends" + fmt.Sprint(AuthorID),
		Value: &result, // destination
		TTL:   time.Hour,
		Do: func(*cache.Item) (interface{}, error) {

			client := http.Client{}
			resp, err := client.Get(`https://only-one.su/api/friends/array_friends/` + fmt.Sprint(AuthorID))
			if err != nil {
				log.Fatalln(err)
			}

			var result []int
			json.NewDecoder(resp.Body).Decode(&result)

			return result, err
		},
	})
	return result, err

}
