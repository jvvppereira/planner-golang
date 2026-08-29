package api

import (
	"context"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
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
