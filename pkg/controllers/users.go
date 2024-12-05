package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"practice-go/pkg/config"
	"practice-go/pkg/jwt"
	"practice-go/pkg/models"

	"github.com/google/uuid"

	"practice-go/pkg/utils"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"
)

func Signup(w http.ResponseWriter, r *http.Request) {

	db := config.ConnectDB()

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Hash the password before saving it to the database
	hashedPassword, err := jwt.HashPassword(user.Password)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, "Could not hash password")
		return
	}

	// Insert the user into the database
	_, err = db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", user.Name, user.Email, hashedPassword)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessHandler(w, r, http.StatusCreated, user)
}

// Login handler for authenticating users and issuing a JWT
func Login(w http.ResponseWriter, r *http.Request) {

	// Initialize database connection
	db := config.ConnectDB()

	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	// Find the user by email
	var storedPassword string
	err = db.QueryRow("SELECT password FROM users WHERE email = ?", user.Email).Scan(&storedPassword)

	fmt.Println(user)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Check if the password matches
	if !jwt.CheckPasswordHash(user.Password, storedPassword) {
		utils.ErrorHandler(w, r, http.StatusUnauthorized, "Invalid email or password")
		return
	}

	// Generate a JWT token
	token, err := jwt.GenerateJWT(user.Email)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, "Could not generate token")
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",                 // Cookie name
		Value:    token,                   // JWT token value
		Path:     "/",                     // Cookie will be sent for all paths
		HttpOnly: true,                    // Prevent client-side access to the cookie (for security)
		Secure:   false,                   // Use secure cookies in HTTPS (set to true for production)
		SameSite: http.SameSiteStrictMode, // Prevent CSRF attacks
	})
	utils.SuccessHandler(w, r, http.StatusOK, "Login successful")
}

func GetUsers(w http.ResponseWriter, r *http.Request) {

	db := config.ConnectDB()

	var users []models.User
	rows, err := db.Query("SELECT id, name, email, created_at, password FROM users")

	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		var id string
		if err := rows.Scan(&id, &user.Name, &user.Email, &user.CreatedAt, &user.Password); err != nil {
			utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		user.ID, _ = uuid.Parse(id)
		users = append(users, user)
	}
	utils.SuccessHandler(w, r, http.StatusOK, users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {

	db := config.ConnectDB()

	id := mux.Vars(r)["id"]

	var user models.User
	var uuidID string

	err := db.QueryRow("SELECT id, name, email, created_at FROM users WHERE id = ?", id).
		Scan(&uuidID, &user.Name, &user.Email, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorHandler(w, r, http.StatusNotFound, "User not found")
		} else {
			utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	user.ID, _ = uuid.Parse(uuidID)
	utils.SuccessHandler(w, r, http.StatusOK, user)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {

	db := config.ConnectDB()

	id := mux.Vars(r)["id"]

	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	_, err = db.Exec("UPDATE users SET name = ?, email = ? WHERE id = ?",
		user.Name, user.Email, id)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	user.ID, _ = uuid.Parse(id)
	utils.SuccessHandler(w, r, http.StatusOK, user)
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {

	db := config.ConnectDB()

	id := mux.Vars(r)["id"]

	_, err := db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	utils.SuccessHandler(w, r, http.StatusNoContent, "User deleted successfully")
}
