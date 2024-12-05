package middleware

import (
	"fmt"
	"net/http"
	"practice-go/pkg/utils"

	"github.com/dgrijalva/jwt-go"
)

func Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Retrieve the token from the cookie
		cookie, err := r.Cookie("token")
		if err != nil {
			utils.ErrorHandler(w, r, http.StatusUnauthorized, "Unauthorized: No token provided")
			return
		}

		// Parse and validate the token
		tokenString := cookie.Value
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Validate the signing method
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte("your_secret_key"), nil // Same secret key as in GenerateJWT
		})

		if err != nil || !token.Valid {
			utils.ErrorHandler(w, r, http.StatusUnauthorized, "Unauthorized: Invalid token")
			return
		}
		next.ServeHTTP(w, r)
	})
}
