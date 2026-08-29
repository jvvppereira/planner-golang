package api

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"

	"planner-golang/internal/api/spec"
	"planner-golang/internal/pgstore"
)

const (
	errInvalidUUID         = "invalid uuid"
	errSomethingWentWrong  = "something went wrong, try again"
	errInvalidJSON         = "invalid JSON: "
	errInvalidInput        = "invalid input: "
	errTripNotFound        = "trip not found"
	errParticipantNotFound = "participant not found"
	contentTypeJSON        = "application/json"
)

type store interface {
	CreateTrip(context.Context, pgxpool.Pool, spec.PostTripsJSONBody) (uuid.UUID, error)
	GetParticipant(ctx context.Context, participantID uuid.UUID) (pgstore.Participant, error)
	ConfirmParticipant(ctx context.Context, participantID uuid.UUID) error
	GetTrip(ctx context.Context, id uuid.UUID) (pgstore.Trip, error)
	UpdateTrip(ctx context.Context, arg pgstore.UpdateTripParams) error
	GetTripActivites(ctx context.Context, tripID uuid.UUID) ([]pgstore.Activity, error)
	CreateActivity(ctx context.Context, arg pgstore.CreateActivityParams) (uuid.UUID, error)
	InvitePaticipantToTrip(ctx context.Context, arg pgstore.InvitePaticipantToTripParams) (uuid.UUID, error)
	GetTripLinks(ctx context.Context, tripID uuid.UUID) ([]pgstore.Link, error)
	CreatTripLink(ctx context.Context, arg pgstore.CreatTripLinkParams) (uuid.UUID, error)
	GetParticipants(ctx context.Context, tripID uuid.UUID) ([]pgstore.Participant, error)
}

type mailer interface {
	SendConfirmTripEmailToTripOwner(tripId uuid.UUID) error
}

type API struct {
	store     store
	logger    *zap.Logger
	validator *validator.Validate
	pool      *pgxpool.Pool
	mailer    mailer
}

func NewAPI(pool *pgxpool.Pool, logger *zap.Logger, mailer mailer) API {
	validator := validator.New(validator.WithRequiredStructEnabled())
	return API{pgstore.New(pool), logger, validator, pool, mailer}
}

// Confirms a participant on a trip.
// (PATCH /participants/{participantId}/confirm)
func (api API) PatchParticipantsParticipantIDConfirm(w http.ResponseWriter, r *http.Request, participantID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse((participantID))
	if err != nil {
		resp = spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		participant, err := api.store.GetParticipant(r.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				resp = spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: errParticipantNotFound})
			} else {
				api.logger.Error("falied to get the participant", zap.Error(err), zap.String("paticipant_id", participantID))
				resp = spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: errSomethingWentWrong})
			}
		} else if participant.IsConfirmed {
			resp = spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: "participant already confirmed"})
		} else if err := api.store.ConfirmParticipant(r.Context(), id); err != nil {
			api.logger.Error("falied to confirm participant", zap.Error(err), zap.String("paticipant_id", participantID))
			resp = spec.PatchParticipantsParticipantIDConfirmJSON400Response(spec.Error{Message: errSomethingWentWrong})
		} else {
			resp = spec.PatchParticipantsParticipantIDConfirmJSON204Response(nil)
		}
	}

	return resp
}

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

// Get a trip details.
// (GET /trips/{tripId})
func (api API) GetTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse(tripID)
	if err != nil {
		resp = spec.GetTripsTripIDJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		trip, err := api.store.GetTrip(r.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				resp = spec.GetTripsTripIDJSON400Response(spec.Error{Message: errTripNotFound})
			} else {
				api.logger.Error("failed to get trip", zap.Error(err), zap.String("trip_id", tripID))
				resp = spec.GetTripsTripIDJSON400Response(spec.Error{Message: errSomethingWentWrong})
			}
		} else {
			resp = spec.GetTripsTripIDJSON200Response(struct {
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
	}

	return resp
}

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

// Get a trip activities.
// (GET /trips/{tripId}/activities)
func (api API) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *spec.Response {
	var resp *spec.Response

	id, err := uuid.Parse(tripID)
	if err != nil {
		resp = spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: errInvalidUUID})
	} else {
		activities, err := api.store.GetTripActivites(r.Context(), id)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				resp = spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: errTripNotFound})
			} else {
				api.logger.Error("failed to get trip activities", zap.Error(err), zap.String("trip_id", tripID))
				resp = spec.GetTripsTripIDActivitiesJSON400Response(spec.Error{Message: errSomethingWentWrong})
			}
		} else {
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
			resp = spec.GetTripsTripIDActivitiesJSON200Response(result)
		}
	}

	return resp
}

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
			resp = spec.GetTripsTripIDLinksJSON200Response(result)
		}
	}

	return resp
}

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
			resp = spec.GetTripsTripIDParticipantsJSON200Response(result)
		}
	}

	return resp
}
