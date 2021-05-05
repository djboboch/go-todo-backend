package posts

import (
	"encoding/json"
	"fmt"
	"github.com/djboboch/go-todo/handlers"
	"github.com/djboboch/go-todo/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strings"
)

type CreatePostItemRequest struct {
	Content string `json:"content"`
}

func Get(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var posts []models.Post

		posts, err = models.AllPosts(env.DB)
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
}

func Create(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		if r.Header.Get("Content-Type") == "application/json" {
			var createPostRequest CreatePostItemRequest

			err = json.NewDecoder(r.Body).Decode(&createPostRequest)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(strings.Join([]string{"500 Internal server error -", err.Error()}, " ")))
			}

			err = models.CreatePost(env.DB, createPostRequest.Content)
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
}

func Delete(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var err error

		vars := mux.Vars(r)

		err = models.DeletePost(env.DB, vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Post Deleted"))
	}
}
