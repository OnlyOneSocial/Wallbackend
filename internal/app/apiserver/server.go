package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/katelinlis/Wallbackend/internal/app/store"
	"github.com/sirupsen/logrus"
)

type server struct {
	router    *mux.Router
	logger    *logrus.Logger
	store     store.Store
	HTTPstore store.HTTPStore
	redis     *redis.Client
	cache     *cache.Cache
}

const (
	ctxKeyUser ctxKey = iota
)

type ctxKey int8

var (
	errIncorrectEmailOrPassword = errors.New("incorect email or password")
	jwtsignkey                  string
)

func newServer(store store.Store, httpstore store.HTTPStore, config *Config) *server {
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	mycache := cache.New(&cache.Options{
		Redis:      rdb,
		LocalCache: cache.NewTinyLFU(10000, time.Minute),
	})

	s := &server{
		router:    mux.NewRouter(),
		logger:    logrus.New(),
		cache:     mycache,
		redis:     rdb,
		store:     store,
		HTTPstore: httpstore,
	}
	s.configureRouter()

	jwtsignkey = config.JwtSignKey

	return s
}

type UserAuth struct {
	LegacyID int
	UserUUID string
	Username string
}

func verifyTokenRS256(token, publicKeyPath string) (userAuth UserAuth, err error) {
	keyData, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return userAuth, err
	}
	key, err := jwt.ParseRSAPublicKeyFromPEM(keyData)
	if err != nil {
		return userAuth, err
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			msg := fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
			return 0, msg
		}
		return key, nil
	})

	if parsedToken != nil && parsedToken.Valid {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			username := claims["preferred_username"].(string)
			UserUUID := claims["sub"].(string)

			return UserAuth{
				UserUUID: UserUUID,
				Username: username,
			}, nil
		}
	}

	return userAuth, nil
}

func (s *server) GetDataFromToken(w http.ResponseWriter, r *http.Request) (userAuthData UserAuth, err error) {
	var token string
	tokens, ok := r.Header["Authorization"]
	if ok && len(tokens) >= 1 {
		token = tokens[0]
		token = strings.TrimPrefix(token, "Bearer ")
	}

	if token == "" {
		return userAuthData, errors.New("Token is missing")
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			msg := fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			return 0, msg
		}
		return []byte(jwtsignkey), nil
	})

	if parsedToken == nil || parsedToken.Method == nil {
		return userAuthData, errors.New("Error parsing token")
	}

	if parsedToken.Method.Alg() != "HS256" {
		path, err := os.Getwd()
		fmt.Println(path)
		if err != nil {
			log.Fatal(err)
		}
		userAuthData, err := verifyTokenRS256(token, path+"/configs/public_key.key")
		if err != nil {
			return userAuthData, errors.New("Error parsing token")
		}

		userid, err := s.HTTPstore.User().GetUser(0, userAuthData.UserUUID)
		if err != nil {
			return userAuthData, errors.New("error find legacy userid by ssoid " + err.Error())
		}
		userAuthData.LegacyID = userid.ID

		return userAuthData, err

	}

	if err != nil {
		//s.error(w, r, http.StatusUnauthorized, errors.New("Error parsing token"))
		return userAuthData, errors.New("Error parsing token")
	}
	if parsedToken != nil && parsedToken.Valid {
		if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok {
			userid := claims["userid"].(float64)
			userAuthData.LegacyID = int(userid)
			return userAuthData, nil
		}
	}
	return userAuthData, nil

}

func (s *server) UrlLimitOffset(request *http.Request) (int, int) {
	var offset int
	var limit = 20
	urlParams := request.URL.Query()
	if len(urlParams["offset"]) > 0 {
		offset_, err := strconv.Atoi(urlParams["offset"][0])
		if err == nil {
			offset = offset_
		}
	}
	if len(urlParams["limit"]) > 0 {
		limit_, err := strconv.Atoi(urlParams["limit"][0])
		if err == nil {
			limit = limit_
		}
	}
	return offset, limit
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")

	defer r.Body.Close()
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {

	s.router.Use(s.loggingMiddleware)

	s.router.Methods("OPTIONS").HandlerFunc(
		func(rw http.ResponseWriter, r *http.Request) {
			rw.Header().Set("Access-Control-Allow-Origin", "*")
			rw.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS")
			rw.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
			rw.WriteHeader(http.StatusOK)
		})

	s.ConfigureWallRouter()
}

func (s *server) emptyresponse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.logger.Debug("Request " + r.RequestURI + " from " + r.RemoteAddr)

		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)

	r.Body.Close()

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
