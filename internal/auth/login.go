package auth

import (
	"backend/api"
	"backend/internal/tools"
	"backend/internal/utils"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Request struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	if r.Method != http.MethodPost {
		api.RequestErrorHandler(w, errors.New("invalid method"))
		return
	}

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid json body"))
		return
	}

	user, err := (*db).LoginUser(req.Email, req.Password)

	if err != nil {
		log.Error(err)
		err = errors.New("invalid credentials")
		api.RequestErrorHandler(w, err)
		return
	}

	accessToken, err := utils.CreateAccessToken(user.UserID)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   900,
	})

	if err := json.NewEncoder(w).Encode(api.Response[api.Message]{
		Success: true,
		Code:    http.StatusOK,
		Data:    api.Message{Message: "Login successful"},
	}); err != nil {
		log.Error(err)
	}

}
