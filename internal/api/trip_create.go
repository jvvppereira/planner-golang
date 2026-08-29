package api

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"planner-golang/internal/api/spec"
)

// Create a new trip
// (POST /trips)
func (api API) PostTrips(w http.ResponseWriter, r *http.Request) *spec.Response {
	var body spec.PostTripsJSONBody
	err := json.NewDecoder(r.Body).Decode(&body)

	if err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: errInvalidJSON + err.Error()})
	}

	if err := api.validator.Struct(body); err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: errInvalidInput + err.Error()})
	}

	tripID, err := api.store.CreateTrip(r.Context(), *api.pool, body)
	if err != nil {
		return spec.PostTripsJSON400Response(spec.Error{Message: "falied to create trip, try again"})
	}

	go func() {
		if err := api.mailer.SendConfirmTripEmailToTripOwner(tripID); err != nil {
			api.logger.Error("falied to send email on PostTrips", zap.Error(err), zap.String("tripID", tripID.String()))
		}
	}()

	return spec.PostTripsJSON201Response(spec.CreateTripResponse{TripID: tripID.String()})
}
