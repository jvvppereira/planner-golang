package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"planner-golang/internal/api/spec"
)

// Get a trip links.
// (GET /trips/{tripId}/links)
func (api API) GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse(tripID)
	if err != nil {
		resp = spec.GetTripsTripIDLinksJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		links, err := api.store.GetTripLinks(r.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				resp = spec.GetTripsTripIDLinksJSON400Response(spec.Error{Message: errTripNotFound})
			} else {
				api.logger.Error("failed to get trip links", zap.Error(err), zap.String("trip_id", tripID))
				resp = spec.GetTripsTripIDLinksJSON400Response(spec.Error{Message: errSomethingWentWrong})
			}
		} else {
			resp = api.buildLinksResponse(links)
		}
	}

	return resp
}
