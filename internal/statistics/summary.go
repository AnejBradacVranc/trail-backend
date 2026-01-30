package statistics

import (
	"backend/api"
	"backend/internal/tools"
	"encoding/json"
	"errors"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func GetSummary(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	userEmail := r.URL.Query().Get("user_email")

	if userEmail == "" {
		log.Error("User email is empty")
		api.RequestErrorHandler(w, errors.New("user email cannot be empty"))
		return
	}

	statuses, err := (*db).GetStatisticsSummary(userEmail)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(api.Response[[]*tools.StatisticsSummary]{
		Success: true,
		Code:    http.StatusOK,
		Data:    statuses,
	}); err != nil {
		log.Error(err)
	}

}
