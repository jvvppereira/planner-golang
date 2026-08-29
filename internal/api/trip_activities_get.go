package api

import (
	"net/http"

	"planner-golang/internal/api/spec"
)

// Get a trip activities.
// (GET /trips/{tripId}/activities)
func (api API) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	return api.handleGetTripCollection(
		r,
		tripID,
		api.store.GetTripActivites,
		spec.GetTripsTripIDActivitiesJSON400Response,
		api.buildActivitiesResponse,
	)
}
