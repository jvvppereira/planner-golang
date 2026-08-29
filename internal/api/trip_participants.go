package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"planner-golang/internal/api/spec"
)

// Get a trip participants.
// (GET /trips/{tripId}/participants)
func (api API) GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse(tripID)
	if err != nil {
		resp = spec.GetTripsTripIDParticipantsJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		participants, err := api.store.GetParticipants(r.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				resp = spec.GetTripsTripIDParticipantsJSON400Response(spec.Error{Message: errTripNotFound})
			} else {
				api.logger.Error("failed to get trip participants", zap.Error(err), zap.String("trip_id", tripID))
				resp = spec.GetTripsTripIDParticipantsJSON400Response(spec.Error{Message: errSomethingWentWrong})
			}
		} else {
			resp = api.buildParticipantsResponse(participants)
		}
	}

	return resp
}
