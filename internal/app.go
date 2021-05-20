package internal

import (
	"database/sql"
	"github.com/djboboch/go-todo/handlers/posts"
	"github.com/djboboch/go-todo/internal/http/middleware"
	"github.com/djboboch/go-todo/models"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
)

var (
	host     = os.Getenv("DB_HOST")
	port     = os.Getenv("DB_PORT")
	user     = os.Getenv("DB_USER")
	password = os.Getenv("DB_PASSWORD")
	dbname   = os.Getenv("DB_NAME")
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize() {

	var err error

	a.DB, err = createConnectionPool()
	a.Router = createRouter()

	if err != nil {
		panic(err)
	}

	return
}

func (a *App) CreateRoutes() {

	env := &posts.Env{
		Posts: models.PostModel{
			DB: a.DB,
		},
	}

	a.Router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	apiV1 := a.Router.PathPrefix("/api/v1").Subrouter()

	apiV1.HandleFunc("/todo", env.GetPosts()).Methods(http.MethodGet)

	apiV1.HandleFunc("/todo", env.CreatePost()).Methods(http.MethodPost)

	apiV1.HandleFunc("/todo", env.UpdatePost()).Methods(http.MethodPut)

	apiV1.HandleFunc("/todo/{id}", env.DeletePost()).Methods(http.MethodDelete)

}

func (a *App) AddMiddleware() {
	a.Router.Use(middleware.SetJSONContentType)
}

func (a *App) Run() {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	})

	log.Fatal(http.ListenAndServe(":8000", c.Handler(a.Router)))
}
