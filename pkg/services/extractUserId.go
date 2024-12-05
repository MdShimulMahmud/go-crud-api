package services

import (
	"net/http"
	"practice-go/pkg/config"
	"practice-go/pkg/jwt"
	"practice-go/pkg/utils"

	"github.com/google/uuid"
)

func ExtractUserID(w http.ResponseWriter, r *http.Request) uuid.UUID {
	db := config.ConnectDB()
	cookie, err := r.Cookie("token")

	if err != nil {
		utils.ErrorHandler(w, r, http.StatusUnauthorized, "Unauthorized: No token found")
		return uuid.Nil
	}

	userEmail, err := jwt.ExtractUserEmailFromToken(cookie.Value)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusUnauthorized, "Unauthorized: Invalid token")
		return uuid.Nil
	}

	var userID uuid.UUID

	err = db.QueryRow("SELECT id FROM users WHERE email = ?", userEmail).Scan(&userID)
	if err != nil {
		utils.ErrorHandler(w, r, http.StatusInternalServerError, err.Error())
		return uuid.Nil
	}

	return userID

}
