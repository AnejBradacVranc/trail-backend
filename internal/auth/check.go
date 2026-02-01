package auth

import (
	"backend/api"
	"backend/internal/tools"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func Check(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	userIDNum, err := strconv.ParseInt(userID, 10, 64)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
	}

	user, err := (*db).GetUserByID(userIDNum)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if err := json.NewEncoder(w).Encode(api.Response[tools.UserPublic]{
		Success: true,
		Code:    http.StatusOK,
		Data:    tools.UserPublic{UserID: user.UserID, Name: user.Name, Surname: user.Surname, Email: user.Email, CreatedAt: user.CreatedAt},
	}); err != nil {
		log.Error(err)
	}
}
