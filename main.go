package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/djboboch/go-todo/models"
	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
)

type Env struct {
	db *sql.DB
}

type CreatePostItemRequest struct {
	Content string `json:"content"`
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	for retryCount := 30; retryCount >= 0; retryCount-- {
		err = db.Ping()
		if err != nil {
			if retryCount == 0 {
				log.Fatalf("Not able to establish a connection to the database at %v", psqlInfo)
			}
			fmt.Println("Could not connect to DB")
			time.Sleep(2 * time.Second)
			continue
		}
		break
	}

	env := &Env{
		db: db,
	}

	fmt.Println("Successfully connected!")

	r := mux.NewRouter()

	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	apiV1.HandleFunc("/todo", env.getPost).Methods(http.MethodGet, http.MethodOptions)

	apiV1.HandleFunc("/todo", env.createPost).Methods(http.MethodPost)

	apiV1.HandleFunc("/todo/{id}", env.todoHandler).Methods(http.MethodDelete, http.MethodPut)

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(CorsMiddleware)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func (env *Env) getPost(w http.ResponseWriter, r *http.Request) {
	var err error
	var posts []models.Post

	posts, err = models.AllPosts(env.db)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
}

func (env *Env) createPost(w http.ResponseWriter, r *http.Request) {
	var err error

	if r.Header.Get("Content-Type") == "application/json" {
		var createPostRequest CreatePostItemRequest

		err = json.NewDecoder(r.Body).Decode(&createPostRequest)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(strings.Join([]string{"500 Internal server error -", err.Error()}, " ")))
		}

		err = models.CreatePost(env.db, createPostRequest.Content)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		fmt.Printf("Created new post with content: %+v into DB", createPostRequest.Content)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post Created"))

	} else {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte("415 - unsupported media type. Please send JSON"))
	}
}

func (env *Env) todoHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	switch r.Method {
	case http.MethodDelete:

		vars := mux.Vars(r)

		err = models.DeletePost(env.db, vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post Deleted"))
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
