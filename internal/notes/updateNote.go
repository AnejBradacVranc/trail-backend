package notes

import (
	"backend/api"
	"backend/internal/tools"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

type UpdateRequest struct {
	NoteContent string `json:"note_content"`
}

func UpdateNote(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {
	idParam := r.PathValue("id")
	if idParam == "" {
		api.RequestErrorHandler(w, errors.New("note id cannot be empty"))
		return
	}
	noteID, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid note id"))
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, errors.New("invalid json body"))
		return
	}
	if req.NoteContent == "" {
		api.RequestErrorHandler(w, errors.New("missing required fields"))
		return
	}

	if err := (*db).UpdateNote(context.TODO(), nil, noteID, req.NoteContent); err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	api.RequestSuccessHandler(w, "Note updated")
}
