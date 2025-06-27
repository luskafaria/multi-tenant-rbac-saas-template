package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"

	db "github.com/lucasfaria/rbac/api-go/internal/database/sqlc"
	"github.com/lucasfaria/rbac/api-go/internal/middleware"
)

// GetProfile godoc
// @Summary Get authenticated user profile
// @Description Get the profile of the currently authenticated user.
// @Tags auth
// @Accept  json
// @Produce  json
// @Success 200 {object} object{user=object{id=string,name=string,email=string,avatarUrl=string}}
// @Failure 400 {string} string "User not found"
// @Router /profile [get]
func GetProfile(queries *db.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get userId from authentication middleware
		userId, err := middleware.GetCurrentUserID(r)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := queries.GetUser(r.Context(), userId)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "User not found", http.StatusBadRequest)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		response := map[string]interface{}{
			"user": map[string]interface{}{
				"id":        user.ID,
				"name":      user.Name.String,
				"email":     user.Email,
				"avatarUrl": user.AvatarUrl.String,
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(response)
	}
}
