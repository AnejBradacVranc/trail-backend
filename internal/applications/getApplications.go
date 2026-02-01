package applications

import (
	"backend/api"
	"backend/internal/tools"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func GetApplications(w http.ResponseWriter, r *http.Request, db *tools.DbInterface) {

	userID, ok := r.Context().Value("user_id").(string)

	statusIDsStr := r.URL.Query()["status_id"] // returns []string
	var statusIDs []int64

	for _, idStr := range statusIDsStr {
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			api.RequestErrorHandler(w, err)
			continue
		}
		statusIDs = append(statusIDs, id)
	}

	if !ok {
		api.RequestUnauthorisedHandler(w, "Unauthorized")
		return
	}

	userIDNum, err := strconv.ParseInt(userID, 10, 64)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
	}

	applications, err := (*db).GetApplicationsFromUserByID(userIDNum, statusIDs...)

	if err != nil {
		log.Error(err)
		api.RequestErrorHandler(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(api.Response[[]*tools.ApplicationSummary]{
		Success: true,
		Code:    http.StatusOK,
		Data:    applications,
	}); err != nil {
		log.Error(err)
	}

}
