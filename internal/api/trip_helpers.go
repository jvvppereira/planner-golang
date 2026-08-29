package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"planner-golang/internal/api/spec"
	"planner-golang/internal/pgstore"
)

func (api API) handleTripNotFound(err error, tripID string, responseFunc func(spec.Error) *spec.Response) *spec.Response {
	if errors.Is(err, pgx.ErrNoRows) {
		return responseFunc(spec.Error{Message: errTripNotFound})
	}
	api.logger.Error("failed to get trip", zap.Error(err), zap.String("trip_id", tripID))
	return responseFunc(spec.Error{Message: errSomethingWentWrong})
}

func (api API) parseTripID(tripID string) (uuid.UUID, *spec.Response) {
	id, err := uuid.Parse(tripID)
	if err != nil {
		return uuid.UUID{}, spec.PostTripsTripIDActivitiesJSON400Response(spec.Error{Message: errInvalidUUID})
	}
	return id, nil
}

func (api API) decodeAndValidateBody[T any](r *http.Request, responseFunc func(spec.Error) *spec.Response) (*T, *spec.Response) {
	var body T
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, responseFunc(spec.Error{Message: errInvalidJSON + err.Error()})
	}
	if err := api.validator.Struct(body); err != nil {
		return nil, responseFunc(spec.Error{Message: errInvalidInput + err.Error()})
	}
	return &body, nil
}

func (api API) confirmTripAndSendEmail(ctx context.Context, tripID uuid.UUID, trip pgstore.Trip) *spec.Response {
	err := api.store.UpdateTrip(ctx, pgstore.UpdateTripParams{
		Destination: trip.Destination,
		EndsAt:      trip.EndsAt,
		StartsAt:    trip.StartsAt,
		IsConfirmed: true,
		ID:          tripID,
	})
	if err != nil {
		api.logger.Error("failed to confirm trip", zap.Error(err), zap.String("trip_id", tripID.String()))
		return spec.GetTripsTripIDConfirmJSON400Response(spec.Error{Message: errSomethingWentWrong})
	}

	go func() {
		if err := api.mailer.SendConfirmTripEmailToTripOwner(tripID); err != nil {
			api.logger.Error("failed to send email on GetTripsTripIDConfirm", zap.Error(err), zap.String("tripID", tripID.String()))
		}
	}()

	return spec.GetTripsTripIDConfirmJSON204Response(nil)
}

func (api API) buildTripResponse(trip pgstore.Trip) *spec.Response {
	return spec.GetTripsTripIDJSON200Response(struct {
		Trip struct {
			Destination string    `json:"destination"`
			EndsAt      time.Time `json:"ends_at"`
			ID          string    `json:"id"`
			IsConfirmed bool      `json:"is_confirmed"`
			StartsAt    time.Time `json:"starts_at"`
		} `json:"trip"`
	}{
		Trip: struct {
			Destination string    `json:"destination"`
			EndsAt      time.Time `json:"ends_at"`
			ID          string    `json:"id"`
			IsConfirmed bool      `json:"is_confirmed"`
			StartsAt    time.Time `json:"starts_at"`
		}{
			Destination: trip.Destination,
			EndsAt:      trip.EndsAt.Time,
			ID:          trip.ID.String(),
			IsConfirmed: trip.IsConfirmed,
			StartsAt:    trip.StartsAt.Time,
		},
	})
}

func (api API) buildActivitiesResponse(activities []pgstore.Activity) *spec.Response {
	var result spec.TripActivities
	for _, a := range activities {
		result.Activities = append(result.Activities, spec.Activities{
			Date: a.OccursAt.Time,
			Activities: []spec.Activity{{
				ID:       a.ID.String(),
				OccursAt: a.OccursAt.Time,
				Title:    a.Title,
			}},
		})
	}
	return spec.GetTripsTripIDActivitiesJSON200Response(result)
}

func (api API) buildLinksResponse(links []pgstore.Link) *spec.Response {
	var result struct {
		Links []struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			URL   string `json:"url"`
		} `json:"links"`
	}
	for _, l := range links {
		result.Links = append(result.Links, struct {
			ID    string `json:"id"`
			Title string `json:"title"`
			URL   string `json:"url"`
		}{
			ID:    l.ID.String(),
			Title: l.Title,
			URL:   l.Url,
		})
	}
	return spec.GetTripsTripIDLinksJSON200Response(result)
}

func (api API) buildParticipantsResponse(participants []pgstore.Participant) *spec.Response {
	var result struct {
		Participants []struct {
			Email       openapi_types.Email `json:"email"`
			ID          string              `json:"id"`
			IsConfirmed bool                `json:"is_confirmed"`
			Name        *string             `json:"name"`
		} `json:"participants"`
	}
	for _, p := range participants {
		result.Participants = append(result.Participants, struct {
			Email       openapi_types.Email `json:"email"`
			ID          string              `json:"id"`
			IsConfirmed bool                `json:"is_confirmed"`
			Name        *string             `json:"name"`
		}{
			Email:       openapi_types.Email(p.Email),
			ID:          p.ID.String(),
			IsConfirmed: p.IsConfirmed,
			Name:        nil,
		})
	}
	return spec.GetTripsTripIDParticipantsJSON200Response(result)
}
