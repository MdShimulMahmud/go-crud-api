package routes

import (
	"net/http"
	controllers "practice-go/pkg/controllers/books"
	"practice-go/pkg/middleware"

	"github.com/gorilla/mux"
)

func BookRoutes(router *mux.Router) {
	b := router.PathPrefix("/books").Subrouter()

	b.Use(middleware.Authenticate)

	b.HandleFunc("/", controllers.GetBooks).Methods(http.MethodGet) // Public route

	b.HandleFunc("/{id}", controllers.GetBook).Methods(http.MethodGet) // Public route

	b.Handle("/", middleware.Authenticate(http.HandlerFunc(controllers.CreateBook))).Methods(http.MethodPost) // Protected route

	b.Handle("/{id}", middleware.Authenticate(http.HandlerFunc(controllers.UpdateBook))).Methods(http.MethodPut) // Protected route

	b.Handle("/{id}", middleware.Authenticate(http.HandlerFunc(controllers.DeleteBook))).Methods(http.MethodDelete) // Protected route

	b.HandleFunc("/{id}/reviews", controllers.GetBookReviews).Methods(http.MethodGet) // Public route

	b.Handle("/{id}/reviews", middleware.Authenticate(http.HandlerFunc(controllers.CreateBookReview))).Methods(http.MethodPost) // Protected route

}
