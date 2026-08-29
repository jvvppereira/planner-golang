package api

import (
	"net/http"

	"github.com/google/uuid"

	"planner-golang/internal/api/spec"
)

// Confirm a trip and send e-mail invitations.
// (GET /trips/{tripId}/confirm)
func (api API) GetTripsTripIDConfirm(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: errInvalidUUID})
	}

	trip, err := api.store.GetTrip(r.Context(), id)
	if err != nil {
		return api.handleTripNotFound(err, tripID, spec.GetTripsTripIDConfirmJSON400Response)
	} else if trip.IsConfirmed {
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: "trip already confirmed"})
	}

	return api.confirmTripAndSendEmail(r.Context(), id, trip)
}
