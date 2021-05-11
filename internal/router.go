package internal

import "github.com/gorilla/mux"

func createRouter() *mux.Router {
	return mux.NewRouter()
}
