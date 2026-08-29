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

// Create a trip activity.
// (POST /trips/{tripId}/activities)
func (api API) PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse(tripID)
	if err != nil {
		resp = spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		var body spec.PostTripsTripIDActivitiesJSONBody
		err = json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			resp = spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: errInvalidJSON + err.Error()})
		} else if err := api.validator.Struct(body); err != nil {
			resp = spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: errInvalidInput + err.Error()})
		} else {
			activityID, err := api.store.CreateActivity(r.Context(), pgstore.CreateActivityParams{
				TripID:   id,
				Title:    body.Title,
				OccursAt: pgtype.Timestamp{Valid: true, Time: body.OccursAt},
			})
			if err != nil {
				if errors.Is(err, pgx.ErrNoRows) {
					resp = spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: errTripNotFound})
				} else {
					api.logger.Error("failed to create activity", zap.Error(err), zap.String("trip_id", tripID))
					resp = spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: errSomethingWentWrong})
				}
			} else {
				resp = spec.PostTripsTripIDActivitiesJSON201Response(struct {
					ActivityID string `json:"activityId"`
				}{ActivityID: activityID.String()})
			}
		}
	}

	return resp
}
