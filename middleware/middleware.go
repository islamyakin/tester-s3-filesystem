package middleware

import (
	"encoding/base64"
	"github.com/islamyakin/tester-s3-filesystem/app/service"
	"github.com/islamyakin/tester-s3-filesystem/models"
	"net/http"
	"strings"
)

// middleware.go
func BasicAuthMiddleware(next http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

		// Handle preflight requests (OPTIONS)
		if r.Method == http.MethodOptions {
			return
		}
		// Check if the request contains the Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// If Authorization header is missing, return 401 Unauthorized
			w.Header().Set("WWW-Authenticate", `Basic realm="User and Key Authentication"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Decode the Authorization header to get the username and password
		authValue := strings.TrimPrefix(authHeader, "Basic ")
		decodedAuth, err := base64.StdEncoding.DecodeString(authValue)
		if err != nil {
			http.Error(w, "Failed to decode authorization header", http.StatusInternalServerError)
			return
		}

		auth := string(decodedAuth)
		credentials := strings.Split(auth, ":")
		if len(credentials) != 2 {
			http.Error(w, "Invalid authorization credentials", http.StatusUnauthorized)
			return
		}

		// Check if the username and password match the expected values
		username := credentials[0]
		password := credentials[1]

		var user models.User
		result := service.DbUserAuth.Where("username = ?", username).First(&user)
		if result.Error != nil {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}

		if user.Password != service.HashPassword(password) {
			http.Error(w, "Invalid password", http.StatusUnauthorized)
			return
		}

		// Authentication is successful, call the next handler
		next(w, r)
	}
}
