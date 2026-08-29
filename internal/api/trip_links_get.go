package api

import (
	"net/http"

	"planner-golang/internal/api/spec"
)

// Get a trip links.
// (GET /trips/{tripId}/links)
func (api API) GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	return api.handleGetTripCollection(
		r,
		tripID,
		api.store.GetTripLinks,
		spec.GetTripsTripIDLinksJSON400Response,
		api.buildLinksResponse,
	)
}
