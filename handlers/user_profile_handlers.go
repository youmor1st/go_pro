package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/gorilla/mux"
	"net/http"
	"shop/models"
	"strings"
)

func GetProfileHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)
	if currentUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentUser)
}
func getCurrentUser(r *http.Request) *models.User {
	authorizationHeader := r.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil
	}
	tokenString := strings.Replace(authorizationHeader, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("your_secret_key"), nil
	})
	if err != nil {
		fmt.Println("Error parsing token:", err)
		return nil
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)

		user, err := models.GetUserByUsername(username)
		if err != nil {
			fmt.Println("Error fetching user:", err)
			return nil
		}
		return user
	}
	return nil
}
func UpdateProfileHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)

	if currentUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	var updatedProfile models.User
	err := json.NewDecoder(r.Body).Decode(&updatedProfile)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if currentUser.ID != updatedProfile.ID {
		http.Error(w, "You can only update your own profile", http.StatusForbidden)
		return
	}

	err = models.UpdateUserProfile(&updatedProfile)
	if err != nil {
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func DeleteProfileHandler(w http.ResponseWriter, r *http.Request) {
	currentUser := getCurrentUser(r)

	if currentUser == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	err := models.DeleteUserProfile(currentUser.ID)
	if err != nil {
		http.Error(w, "Failed to delete profile", http.StatusInternalServerError)
		return
	}

	// Отправка успешного ответа
	w.WriteHeader(http.StatusOK)
}

func GetUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]

	user, err := models.GetUserByUsername(username)
	if err != nil {
		http.Error(w, "Failed to get user profile", http.StatusInternalServerError)
		return
	}

	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
