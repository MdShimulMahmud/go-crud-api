package reviews

import (
	"encoding/json"
	"net/http"
	"practice-go/pkg/config"
	"practice-go/pkg/models"
	"practice-go/pkg/utils"

	"github.com/gorilla/mux"
)

func CreateBookReview(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	// Get the logged-in user ID from the request context
	userID := r.Context().Value("userID").(string)

	params := mux.Vars(r)
	bookID := params["bookId"]

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
	bookID := params["bookId"]

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
			utils.ErrorHandler(w, r, http.StatusInternalServerError, "Failed to parse reviews")
			return
		}
		reviews = append(reviews, review)
	}

	utils.SuccessHandler(w, r, http.StatusOK, reviews)
}
