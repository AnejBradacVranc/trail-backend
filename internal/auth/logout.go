package auth

import (
	"backend/api"
	"encoding/json"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
	})

	if err := json.NewEncoder(w).Encode(api.Response[api.Message]{
		Success: true,
		Code:    http.StatusOK,
		Data:    api.Message{Message: "Logout successful"},
	}); err != nil {
		log.Error(err)
	}
}
