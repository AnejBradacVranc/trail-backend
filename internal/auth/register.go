package auth

import (
	"backend/api"
	"backend/internal/tools"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid json body"))
		return
	}

	_, err := (*db).CreateUser(req.Name, req.Surname, req.Email, req.Password)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
	}

	if err := json.NewEncoder(w).Encode(api.Response[api.Message]{
		Success: true,
		Code:    http.StatusOK,
		Data:    api.Message{Message: "Registration successful"},
	}); err != nil {
		log.Error(err)
	}

}
