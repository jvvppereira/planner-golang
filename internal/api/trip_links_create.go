package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"planner-golang/internal/api/spec"
	"planner-golang/internal/pgstore"
)

// Create a trip link.
// (POST /trips/{tripId}/links)
func (api API) PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse(tripID)
	if err != nil {
		resp = spec.PostTripsTripIDLinksJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		var body spec.PostTripsTripIDLinksJSONBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			resp = spec.PostTripsTripIDLinksJSON400Response(spec.Error{Message: errInvalidJSON + err.Error()})
		} else if err := api.validator.Struct(body); err != nil {
			resp = spec.PostTripsTripIDLinksJSON400Response(spec.Error{Message: errInvalidInput + err.Error()})
		} else {
			linkID, err := api.store.CreatTripLink(r.Context(), pgstore.CreatTripLinkParams{
				TripID: id,
				Title:  body.Title,
				Url:    body.URL,
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					resp = spec.PostTripsTripIDLinksJSON400Response(spec.Error{Message: errTripNotFound})
				} else {
					api.logger.Error("failed to create link", zap.Error(err), zap.String("trip_id", tripID))
					resp = spec.PostTripsTripIDLinksJSON400Response(spec.Error{Message: errSomethingWentWrong})
				}
			} else {
				resp = spec.PostTripsTripIDLinksJSON201Response(struct {
					LinkID string `json:"linkId"`
				}{LinkID: linkID.String()})
			}
		}
	}

	return resp
}
