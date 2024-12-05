package main

import (
	"log"
	"net/http"

	"practice-go/pkg/middleware"
	"practice-go/pkg/routes"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	routes.UsersRoutes(router)
	routes.BookRoutes(router)

	routes.BorrowRoutes(router)
	router.Use(middleware.Logger)
	// Start server on port 8000
	log.Fatal(http.ListenAndServe(":8000", router))
}
