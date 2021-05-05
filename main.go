package main

import (
	"database/sql"
	"fmt"
	"github.com/djboboch/go-todo/handlers"
	"github.com/djboboch/go-todo/handlers/posts"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
)

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

	env := &handlers.Env{
		DB: db,
	}

	fmt.Println("Successfully connected!")

	r := mux.NewRouter()

	apiV1 := r.PathPrefix("/api/v1").Subrouter()

	apiV1.HandleFunc("/todo", posts.Get(env)).Methods(http.MethodGet, http.MethodOptions)

	apiV1.HandleFunc("/todo", posts.Create(env)).Methods(http.MethodPost)

	apiV1.HandleFunc("/todo/{id}", posts.Delete(env)).Methods(http.MethodDelete)

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(CorsMiddleware)

	log.Fatal(http.ListenAndServe(":8000", r))
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}
