package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"planner-golang/internal/api/spec"
)

// Get a trip activities.
// (GET /trips/{tripId}/activities)
func (api API) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse(tripID)
	if err != nil {
		resp = spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		activities, err := api.store.GetTripActivites(r.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				resp = spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: errTripNotFound})
			} else {
				api.logger.Error("failed to get trip activities", zap.Error(err), zap.String("trip_id", tripID))
				resp = spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: errSomethingWentWrong})
			}
		} else {
			resp = api.buildActivitiesResponse(activities)
		}
	}

	return resp
}
