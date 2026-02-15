package notes

import (
	"backend/api"
	"backend/internal/tools"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Request struct {
	NoteContent   string `json:"note_content"`
	ApplicationId *int64 `json:"application_id"`
}

func CreateNote(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	var req Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid json body"))
		return
	}

	if req.NoteContent == "" || req.ApplicationId == nil {
		api.RequestErrorHandler(w, errors.New("missing required fields"))
		return
	}

	id, err := (*db).CreateNote(context.TODO(), nil, *req.ApplicationId, req.NoteContent)

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
