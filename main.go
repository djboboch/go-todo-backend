package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/djboboch/go-todo/models"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"strings"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "golang"
	password = "password"
	dbname   = "testdb"
)

type Env struct {
	db *sql.DB
}

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	fmt.Println(psqlInfo)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	env := &Env{
		db: db,
	}

	fmt.Println("Successfully connected!")

	http.HandleFunc("/api/v1/todo", env.todoHandler)

	log.Fatal(http.ListenAndServe(":8000", nil))
}

func (env *Env) todoHandler(w http.ResponseWriter, r *http.Request) {
	var err error

	switch r.Method {
	case http.MethodGet:
		var posts []models.Post

		posts, err = models.AllPosts(env.db)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		err = json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	case http.MethodPost:
		if r.Header.Get("Content-Type") == "application/json" {
			var post models.Post

			err = json.NewDecoder(r.Body).Decode(&post)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(strings.Join([]string{"500 Internal server error -", err.Error()}, " ")))
			}

			err = models.CreatePost(env.db, post)
			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
			}

			fmt.Printf("Inserted %+v into DB", post)

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Post Created"))

		} else {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("415 - unsupported media type. Please send JSON"))
		}

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("405 - Method not allowed"))
	}
}
