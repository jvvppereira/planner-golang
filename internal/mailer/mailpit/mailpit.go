package mailpit

import (
	"context"
	"fmt"
	"planner-golang/internal/pgstore"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/wneessen/go-mail"
)

type store interface {
	GetTrip(context.Context, uuid.UUID) (pgstore.Trip, error)
}

type MailPit struct {
	store store
}

func NewMailPit(pool *pgxpool.Pool) MailPit {
	return MailPit{pgstore.New(pool)}
}

func (mp MailPit) SendConfirmTripEmailToTripOwner(tripId uuid.UUID) error {
	ctx := context.Background()
	var err error

	trip, err := mp.store.GetTrip(ctx, tripId)
	if err == nil {
		msg := mail.NewMsg()
		if err = msg.From("mailpit@planner.com"); err == nil {
			if err = msg.To(trip.OwnerEmail); err == nil {
				msg.Subject("Confirm your trip!")
				msg.SetBodyString(mail.TypeTextPlain, fmt.Sprintf(`
					Hi %s,

					Your trip to %s that starts on %s needs to be confirmed.

					Click on link below to confirm.
				`,
					trip.OwnerName,
					trip.Destination,
					trip.StartsAt.Time.Format(time.DateOnly),
				))

				var client *mail.Client
				client, err = mail.NewClient("mailpit", mail.WithTLSPortPolicy(mail.NoTLS), mail.WithPort(1025))
				if err == nil {
					err = client.DialAndSend(msg)
				}
			}
		}
	}

	if err != nil {
		return fmt.Errorf("mailpit: failed to send confirmation email: %w", err)
	}

	return nil
}
