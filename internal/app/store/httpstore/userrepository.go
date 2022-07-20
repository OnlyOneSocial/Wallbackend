package httpstore

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-redis/cache/v8"
	"github.com/katelinlis/Wallbackend/internal/app/model"
)

//UserRepository ...
type UserRepository struct {
	store *Store
}

//GetUser ...
func (r *UserRepository) GetUser(AuthorID int) (usrObj model.UserObj, err error) {
	/*if val, ok := r.store.CacheUser.Get(strconv.Itoa(AuthorID)); ok {
		return val.(model.UserObj)
	}*/
	userID := strconv.Itoa(AuthorID)

	err = r.store.cache.Once(&cache.Item{
		Key:   "user" + fmt.Sprint(userID),
		Value: usrObj, // destination
		Do: func(*cache.Item) (interface{}, error) {

			client := http.Client{}
			resp, err := client.Get(`https://only-one.su/api/user/get/` + userID)
			if err != nil {
				log.Fatalln(err)
			}

			var result map[string]map[string]string
			json.NewDecoder(resp.Body).Decode(&result)
			usrObj := model.UserObj{
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
		Value: result, // destination
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
