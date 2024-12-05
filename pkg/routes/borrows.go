package routes

import (
	"net/http"
	controllers "practice-go/pkg/controllers/borrows"
	"practice-go/pkg/middleware"

	"github.com/gorilla/mux"
)

func BorrowRoutes(router *mux.Router) {
	b := router.PathPrefix("/books/{id}").Subrouter()

	b.Use(middleware.Authenticate)

	b.Handle("/borrows", middleware.Authenticate(http.HandlerFunc(controllers.CreateBorrow))).Methods(http.MethodPost) // Protected route

}
