package api

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	"planner-golang/internal/api/spec"
	"planner-golang/internal/pgstore"
)

// Update a trip.
// (PUT /trips/{tripId})
func (api API) PutTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse(tripID)
	if err != nil {
		resp = spec.PutTripsTripIDJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		var body spec.PutTripsTripIDJSONBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			resp = spec.PutTripsTripIDJSON400Response(spec.Error{Message: errInvalidJSON + err.Error()})
		} else if err := api.validator.Struct(body); err != nil {
			resp = spec.PutTripsTripIDJSON400Response(spec.Error{Message: errInvalidInput + err.Error()})
		} else {
			err = api.store.UpdateTrip(r.Context(), pgstore.UpdateTripParams{
				Destination: body.Destination,
				EndsAt:      pgtype.Timestamp{Valid: true, Time: body.EndsAt},
				StartsAt:    pgtype.Timestamp{Valid: true, Time: body.StartsAt},
				IsConfirmed: false,
				ID:          id,
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					resp = spec.PutTripsTripIDJSON400Response(spec.Error{Message: errTripNotFound})
				} else {
					api.logger.Error("failed to update trip", zap.Error(err), zap.String("trip_id", tripID))
					resp = spec.PutTripsTripIDJSON400Response(spec.Error{Message: errSomethingWentWrong})
				}
			} else {
				resp = spec.PutTripsTripIDJSON204Response(nil)
			}
		}
	}

	return resp
}
