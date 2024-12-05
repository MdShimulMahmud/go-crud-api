package routes

import (
	"net/http"

	"practice-go/pkg/controllers"
	"practice-go/pkg/middleware"

	"github.com/gorilla/mux"
)

func UsersRoutes(router *mux.Router) {

	u := router.PathPrefix("/users").Subrouter()

	u.HandleFunc("/signup", controllers.Signup).Methods(http.MethodPost) // Signup a new user
	u.HandleFunc("/login", controllers.Login).Methods(http.MethodPost)

	u.Handle("/", middleware.Authenticate(http.HandlerFunc(controllers.GetUsers))).Methods(http.MethodGet)
	u.Handle("/{id}", middleware.Authenticate(http.HandlerFunc(controllers.GetUser))).Methods(http.MethodGet)
	u.Handle("/{id}", middleware.Authenticate(http.HandlerFunc(controllers.UpdateUser))).Methods(http.MethodPut)
	u.Handle("/{id}", middleware.Authenticate(http.HandlerFunc(controllers.DeleteUser))).Methods(http.MethodDelete) // Delete a user by ID

}
