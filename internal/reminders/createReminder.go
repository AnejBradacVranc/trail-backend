package reminders

import (
	"backend/api"
	"backend/internal/tools"
	"backend/internal/utils"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type CreateRequest struct {
	ApplicationID *int64  `json:"application_id"`
	Title         string  `json:"title"`
	Description   *string `json:"description"`
	RemindAt      string  `json:"remind_at"`
	IsCompleted   *bool   `json:"is_completed"`
}

func CreateReminder(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {
	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid json body"))
		return
	}

	if req.ApplicationID == nil || req.Title == "" || req.RemindAt == "" {
		api.RequestErrorHandler(w, errors.New("missing required fields"))
		return
	}

	parsedRemindAt, err := utils.ParseTimeString(req.RemindAt)
	if err != nil {
		api.RequestErrorHandler(w, errors.New("invalid remind_at format"))
		return
	}

	isCompleted := false
	if req.IsCompleted != nil {
		isCompleted = *req.IsCompleted
	}

	id, err := (*db).CreateReminder(*req.ApplicationID, req.Title, req.Description, parsedRemindAt, isCompleted)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(api.Response[api.IDResp]{
		Success: true,
		Code:    http.StatusOK,
		Data:    api.IDResp{ID: id},
	}); err != nil {
		log.Error(err)
	}
}
