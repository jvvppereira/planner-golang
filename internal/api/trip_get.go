package api

import (
	"net/http"

	"github.com/google/uuid"

	"planner-golang/internal/api/spec"
)

// Get a trip details.
// (GET /trips/{tripId})
func (api API) GetTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.GetTripsTripIDJSON400Response(spec.Error{Message: errInvalidUUID})
	}

	trip, err := api.store.GetTrip(r.Context(), id)
	if err != nil {
		return api.handleTripNotFound(err, tripID, spec.GetTripsTripIDJSON400Response)
	}

	return api.buildTripResponse(trip)
}
