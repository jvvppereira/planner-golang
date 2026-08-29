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

// Invite someone to the trip.
// (POST /trips/{tripId}/invites)
func (api API) PostTripsTripIDInvites(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse(tripID)
	if err != nil {
		resp = spec.PostTripsTripIDInvitesJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		var body spec.PostTripsTripIDInvitesJSONBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			resp = spec.PostTripsTripIDInvitesJSON400Response(spec.Error{Message: errInvalidJSON + err.Error()})
		} else if err := api.validator.Struct(body); err != nil {
			resp = spec.PostTripsTripIDInvitesJSON400Response(spec.Error{Message: errInvalidInput + err.Error()})
		} else {
			_, err = api.store.InvitePaticipantToTrip(r.Context(), pgstore.InvitePaticipantToTripParams{
				TripID: id,
				Email:  string(body.Email),
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					resp = spec.PostTripsTripIDInvitesJSON400Response(spec.Error{Message: errTripNotFound})
				} else {
					api.logger.Error("failed to invite participant", zap.Error(err), zap.String("trip_id", tripID))
					resp = spec.PostTripsTripIDInvitesJSON400Response(spec.Error{Message: errSomethingWentWrong})
				}
			} else {
				resp = spec.PostTripsTripIDInvitesJSON201Response(nil)
			}
		}
	}

	return resp
}
