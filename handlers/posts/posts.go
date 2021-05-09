package posts

import (
	"encoding/json"
	"github.com/djboboch/go-todo/internal/http/requests"
	"github.com/djboboch/go-todo/internal/http/responses"
	"github.com/djboboch/go-todo/models"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type Env struct {
	Post models.PostModel
}

func (env *Env) GetPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var posts []models.Post

		posts, err = env.Post.All()

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

func (env *Env) CreatePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		var post *models.Post

		var createPostRequest requests.CreatePostItem

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

		post, err = env.Post.Create(createPostRequest.Content)
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

func (env *Env) UpdatePost() http.HandlerFunc {
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

		err = env.Post.Update(post)
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

func (env *Env) DeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var err error

		vars := mux.Vars(r)

		err = env.Post.Delete(vars["id"])
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
