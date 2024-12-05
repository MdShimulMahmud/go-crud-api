package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"practice-go/pkg/config"
	"practice-go/pkg/jwt"
	"practice-go/pkg/models"
	"practice-go/pkg/services"
	"practice-go/pkg/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// Get all books
func GetBooks(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	rows, err := db.Query("SELECT * FROM books")
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var book models.Book
		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Description, &book.PublishedAt, &book.ImageURL, &book.UserID); err != nil {
			utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
			return
		}
		books = append(books, book)
	}
	utils.SuccessHandler(w, r, http.StatusOK, books)
}

// Get a book by ID
func GetBook(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	params := mux.Vars(r)
	id := params["id"]

	var book models.Book
	err := db.QueryRow("SELECT * FROM books WHERE id = ?", id).Scan(&book.ID, &book.Title, &book.Author, &book.ISBN, &book.Description, &book.PublishedAt, &book.ImageURL, &book.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
		} else {
			utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		}
		return
	}
	utils.SuccessHandler(w, r, http.StatusOK, book)
}

func CreateBook(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()
	cookie, err := r.Cookie("token")

	if err != nil {
		utils.ErrorHandler(w, r, http.StatusUnauthorized, "Unauthorized: No token found")
		return
	}

	userEmail, err := jwt.ExtractUserEmailFromToken(cookie.Value)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusUnauthorized, "Unauthorized: Invalid token")
		return
	}

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	var userID string
	err = db.QueryRow("SELECT id FROM users WHERE email = ?", userEmail).Scan(&userID)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	book.ID = uuid.New()
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, "Invalid user ID")
		return
	}
	book.UserID = userUUID

	_, err = db.Exec(
		"INSERT INTO books (id, title, author, isbn, description, published_at, image_url, user_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)",
		book.ID, book.Title, book.Author, book.ISBN, book.Description, book.PublishedAt, book.ImageURL, book.UserID,
	)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessHandler(w, r, http.StatusCreated, book)
}

// Update a book by ID
func UpdateBook(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	params := mux.Vars(r)
	id := params["id"]

	var book models.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}

	_, err := db.Exec(
		"UPDATE books SET title = ?, author = ?, isbn = ?, description = ?, published_at = ?, image_url = ? WHERE id = ?",
		book.Title, book.Author, book.ISBN, book.Description, book.PublishedAt, book.ImageURL, id)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessHandler(w, r, http.StatusOK, book)
}

// Delete a book by ID
func DeleteBook(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	params := mux.Vars(r)
	id := params["id"]

	_, err := db.Exec("DELETE FROM books WHERE id = ?", id)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessHandler(w, r, http.StatusOK, "Book deleted successfully")
}

func CreateBookReview(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	// Get the logged-in user ID from the request context
	userID := services.ExtractUserID(w, r)
	if userID == uuid.Nil {
		utils.ErrorHandler(w, r, http.StatusUnauthorized, "Unauthorized: Invalid user ID")
		return
	}

	params := mux.Vars(r)
	bookID := params["id"]

	var review models.Review
	if err := json.NewDecoder(r.Body).Decode(&review); err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest, "Invalid input")
		return
	}

	// Insert the review into the database
	_, err := db.Exec(
		"INSERT INTO reviews (id, book_id, user_id, rating, comment) VALUES (UUID(), ?, ?, ?, ?)",
		bookID, userID, review.Rating, review.Comment,
	)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to create review")
		return
	}

	utils.SuccessHandler(w, r, http.StatusCreated, "Review created successfully")
}

func GetBookReviews(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	params := mux.Vars(r)
	bookID := params["id"]

	rows, err := db.Query("SELECT id, user_id, rating, comment, created_at FROM reviews WHERE book_id = ?", bookID)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to fetch reviews")
		return
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var review models.Review
		if err := rows.Scan(&review.ID, &review.UserID, &review.Rating, &review.Comment, &review.CreatedAt); err != nil {
			log.Printf("Error scanning review: %v", err)
			utils.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to parse reviews")
			return
		}
		reviews = append(reviews, review)
	}
	utils.SuccessHandler(w, r, http.StatusOK, reviews)
}
