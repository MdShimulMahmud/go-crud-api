package borrows

import (
	"encoding/json"
	"net/http"
	"practice-go/pkg/config"
	"practice-go/pkg/models"
	"practice-go/pkg/services"
	"practice-go/pkg/utils"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateBorrow(w http.ResponseWriter, r *http.Request) {
	db := config.ConnectDB()
	defer db.Close()

	userID := services.ExtractUserID(w, r)

	var borrow models.Borrow
	if err := json.NewDecoder(r.Body).Decode(&borrow); err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest, err.Error())
		return
	}
	bookId := mux.Vars(r)["id"]
	bookUUID, err := uuid.Parse(bookId)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusBadRequest, "Invalid book ID")
		return
	}

	borrow.BookID = bookUUID
	borrow.UserID = userID

	_, err = db.Exec("INSERT INTO borrows (id, book_id, user_id, status, due_date) VALUES (UUID(),?, ?, ?, ?)", borrow.BookID, borrow.UserID, borrow.Status, borrow.DueDate)

	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return
	}

	utils.SuccessHandler(w, r, http.StatusCreated, borrow)
}

func GetBorrows(w http.ResponseWriter, r *http.Request) {
	// ...
}

func UpdateBorrow(w http.ResponseWriter, r *http.Request) {
	// ...
}

func DeleteBorrow(w http.ResponseWriter, r *http.Request) {
	// ...
}
