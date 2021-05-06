package posts

import (
	"encoding/json"
	"github.com/djboboch/go-todo/handlers"
	"github.com/djboboch/go-todo/models"
	"github.com/djboboch/go-todo/pkg/responses"
	"github.com/gorilla/mux"
	"log"
	"net/http"
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
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.ErrorResponseStatus,
				Content: err.Error(),
			})

			return
		}

		w.WriteHeader(http.StatusOK)
		if len(posts) == 0 {
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.SuccessResponseStatus,
				Content: "You have no todo",
			})
			return
		} else {
			err = json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.SuccessResponseStatus,
				Content: posts,
			})
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				json.NewEncoder(w).Encode(responses.ServerResponse{
					Status:  responses.ErrorResponseStatus,
					Content: err.Error(),
				})

				return
			}
		}
	}
}

func Create(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var post *models.Post

		var createPostRequest CreatePostItemRequest

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.ErrorResponseStatus,
				Content: "Wrong Media type - Send JSON",
			})

			return
		}

		err = json.NewDecoder(r.Body).Decode(&createPostRequest)

		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.ErrorResponseStatus,
				Content: err.Error(),
			})

			return
		}

		post, err = models.CreatePost(env.DB, createPostRequest.Content)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.ErrorResponseStatus,
				Content: err.Error(),
			})

			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(responses.ServerResponse{
			Status:  responses.SuccessResponseStatus,
			Content: &post,
		})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.ErrorResponseStatus,
				Content: err.Error(),
			})

			return
		}
	}
}

func Update(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var post models.Post
		var err error

		if r.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.ErrorResponseStatus,
				Content: "Wrong Media type - Send JSON",
			})

			return
		}

		err = json.NewDecoder(r.Body).Decode(&post)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.ErrorResponseStatus,
				Content: err.Error(),
			})

			return
		}

		err = models.UpdatePost(env.DB, post)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.ErrorResponseStatus,
				Content: err.Error(),
			})

			return
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(responses.ServerResponse{
			Status:  responses.SuccessResponseStatus,
			Content: "Post updated",
		})
	}
}

func Delete(env *handlers.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var err error

		vars := mux.Vars(r)

		err = models.DeletePost(env.DB, vars["id"])
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(responses.ServerResponse{
				Status:  responses.ErrorResponseStatus,
				Content: err.Error(),
			})

			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(responses.ServerResponse{
			Status:  responses.SuccessResponseStatus,
			Content: "Post Deleted",
		})
	}
}
