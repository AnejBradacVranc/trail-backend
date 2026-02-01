package statistics

import (
	"backend/api"
	"backend/internal/tools"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func GetSummary(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	userID, ok := r.Context().Value("user_id").(string)

	if !ok {
		api.RequestUnauthorisedHandler(w, "Unauthorized")
		return
	}

	userIDNum, err := strconv.ParseInt(userID, 10, 64)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
	}

	statuses, err := (*db).GetStatisticsSummary(userIDNum)

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
