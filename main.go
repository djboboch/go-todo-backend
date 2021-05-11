package main

import (
	"github.com/djboboch/go-todo/internal"
	_ "github.com/lib/pq"
)

func main() {

	a := internal.App{}
	a.Initialize()
	a.CreateRoutes()
	a.AddMiddleware()
	a.Run()
}
