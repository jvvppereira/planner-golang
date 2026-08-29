package api

import (
	"net/http"

	"planner-golang/internal/api/spec"
)

// Get a trip participants.
// (GET /trips/{tripId}/participants)
func (api API) GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	return api.handleGetTripCollection(
		r,
		tripID,
		api.store.GetParticipants,
		spec.GetTripsTripIDParticipantsJSON400Response,
		api.buildParticipantsResponse,
	)
}
