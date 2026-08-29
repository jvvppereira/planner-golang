package api

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"

	"planner-golang/internal/api/spec"
)

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
