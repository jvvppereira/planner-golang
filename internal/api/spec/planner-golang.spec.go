// Package spec provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/discord-gophers/goapi-gen version v0.3.0 DO NOT EDIT.
package spec

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/discord-gophers/goapi-gen/runtime"
	openapi_types "github.com/discord-gophers/goapi-gen/types"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
)

// Activities defines model for Activities.
type Activities struct {
	Activities []Activity `json:"activities"`
	Date       time.Time  `json:"date"`
}

// Activity defines model for Activity.
type Activity struct {
	ID       string    `json:"id"`
	OccursAt time.Time `json:"occurs_at"`
	Title    string    `json:"title"`
}

// Bad request
type Error struct {
	Message string `json:"message"`
}

// TripActivities defines model for TripActivities.
type TripActivities struct {
	Activities []Activities `json:"activities"`
}

// PostTripsJSONBody defines parameters for PostTrips.
type PostTripsJSONBody struct {
	Destination    string                `json:"destination" validate:"required,min=4"`
	EmailsToInvite []openapi_types.Email `json:"emails_to_invite" validate:"required,dive,email"`
	EndsAt         time.Time             `json:"ends_at" validate:"required"`
	OwnerEmail     openapi_types.Email   `json:"owner_email" validate:"required,email"`
	OwnerName      string                `json:"owner_name" validate:"required"`
	StartsAt       time.Time             `json:"starts_at" validate:"required"`
}

// PutTripsTripIDJSONBody defines parameters for PutTripsTripID.
type PutTripsTripIDJSONBody struct {
	Destination string    `json:"destination" validate:"required,min=4"`
	EndsAt      time.Time `json:"ends_at" validate:"required"`
	StartsAt    time.Time `json:"starts_at" validate:"required"`
}

// PostTripsTripIDActivitiesJSONBody defines parameters for PostTripsTripIDActivities.
type PostTripsTripIDActivitiesJSONBody struct {
	OccursAt time.Time `json:"occurs_at" validate:"required"`
	Title    string    `json:"title" validate:"required"`
}

// PostTripsTripIDInvitesJSONBody defines parameters for PostTripsTripIDInvites.
type PostTripsTripIDInvitesJSONBody struct {
	Email openapi_types.Email `json:"email" validate:"required,email"`
}

// PostTripsTripIDLinksJSONBody defines parameters for PostTripsTripIDLinks.
type PostTripsTripIDLinksJSONBody struct {
	Title string `json:"title" validate:"required"`
	URL   string `json:"url" validate:"required,url"`
}

// PostTripsJSONRequestBody defines body for PostTrips for application/json ContentType.
type PostTripsJSONRequestBody PostTripsJSONBody

// Bind implements render.Binder.
func (PostTripsJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PutTripsTripIDJSONRequestBody defines body for PutTripsTripID for application/json ContentType.
type PutTripsTripIDJSONRequestBody PutTripsTripIDJSONBody

// Bind implements render.Binder.
func (PutTripsTripIDJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostTripsTripIDActivitiesJSONRequestBody defines body for PostTripsTripIDActivities for application/json ContentType.
type PostTripsTripIDActivitiesJSONRequestBody PostTripsTripIDActivitiesJSONBody

// Bind implements render.Binder.
func (PostTripsTripIDActivitiesJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostTripsTripIDInvitesJSONRequestBody defines body for PostTripsTripIDInvites for application/json ContentType.
type PostTripsTripIDInvitesJSONRequestBody PostTripsTripIDInvitesJSONBody

// Bind implements render.Binder.
func (PostTripsTripIDInvitesJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// PostTripsTripIDLinksJSONRequestBody defines body for PostTripsTripIDLinks for application/json ContentType.
type PostTripsTripIDLinksJSONRequestBody PostTripsTripIDLinksJSONBody

// Bind implements render.Binder.
func (PostTripsTripIDLinksJSONRequestBody) Bind(*http.Request) error {
	return nil
}

// Response is a common response struct for all the API calls.
// A Response object may be instantiated via functions for specific operation responses.
// It may also be instantiated directly, for the purpose of responding with a single status code.
type Response struct {
	body        interface{}
	Code        int
	contentType string
}

// Render implements the render.Renderer interface. It sets the Content-Type header
// and status code based on the response definition.
func (resp *Response) Render(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", resp.contentType)
	render.Status(r, resp.Code)
	return nil
}

// Status is a builder method to override the default status code for a response.
func (resp *Response) Status(code int) *Response {
	resp.Code = code
	return resp
}

// ContentType is a builder method to override the default content type for a response.
func (resp *Response) ContentType(contentType string) *Response {
	resp.contentType = contentType
	return resp
}

// MarshalJSON implements the json.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalJSON() ([]byte, error) {
	return json.Marshal(resp.body)
}

// MarshalXML implements the xml.Marshaler interface.
// This is used to only marshal the body of the response.
func (resp *Response) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return e.Encode(resp.body)
}

// PatchParticipantsParticipantIDConfirmJSON204Response is a constructor method for a PatchParticipantsParticipantIDConfirm response.
// A *Response is returned with the configured status code and content type from the spec.
func PatchParticipantsParticipantIDConfirmJSON204Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        204,
		contentType: "application/json",
	}
}

// PatchParticipantsParticipantIDConfirmJSON400Response is a constructor method for a PatchParticipantsParticipantIDConfirm response.
// A *Response is returned with the configured status code and content type from the spec.
func PatchParticipantsParticipantIDConfirmJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostTripsJSON201Response is a constructor method for a PostTrips response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsJSON201Response(body struct {
	TripID string `json:"tripId"`
}) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}


// CreateTripResponse defines model for CreateTripResponse.
type CreateTripResponse struct {
	TripID string `json:"tripId"`
}

// PostTripsJSON400Response is a constructor method for a PostTrips response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDJSON200Response is a constructor method for a GetTripsTripID response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDJSON200Response(body struct {
	Trip struct {
		Destination string    `json:"destination"`
		EndsAt      time.Time `json:"ends_at"`
		ID          string    `json:"id"`
		IsConfirmed bool      `json:"is_confirmed"`
		StartsAt    time.Time `json:"starts_at"`
	} `json:"trip"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetTripsTripIDJSON400Response is a constructor method for a GetTripsTripID response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PutTripsTripIDJSON204Response is a constructor method for a PutTripsTripID response.
// A *Response is returned with the configured status code and content type from the spec.
func PutTripsTripIDJSON204Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        204,
		contentType: "application/json",
	}
}

// PutTripsTripIDJSON400Response is a constructor method for a PutTripsTripID response.
// A *Response is returned with the configured status code and content type from the spec.
func PutTripsTripIDJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDActivitiesJSON200Response is a constructor method for a GetTripsTripIDActivities response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDActivitiesJSON200Response(body TripActivities) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetTripsTripIDActivitiesJSON400Response is a constructor method for a GetTripsTripIDActivities response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDActivitiesJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostTripsTripIDActivitiesJSON201Response is a constructor method for a PostTripsTripIDActivities response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDActivitiesJSON201Response(body struct {
	ActivityID string `json:"activityId"`
}) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostTripsTripIDActivitiesJSON400Response is a constructor method for a PostTripsTripIDActivities response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDActivitiesJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDConfirmJSON204Response is a constructor method for a GetTripsTripIDConfirm response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDConfirmJSON204Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        204,
		contentType: "application/json",
	}
}

// GetTripsTripIDConfirmJSON400Response is a constructor method for a GetTripsTripIDConfirm response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDConfirmJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostTripsTripIDInvitesJSON201Response is a constructor method for a PostTripsTripIDInvites response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDInvitesJSON201Response(body interface{}) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostTripsTripIDInvitesJSON400Response is a constructor method for a PostTripsTripIDInvites response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDInvitesJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDLinksJSON200Response is a constructor method for a GetTripsTripIDLinks response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDLinksJSON200Response(body struct {
	Links []struct {
		ID    string `json:"id"`
		Title string `json:"title"`
		URL   string `json:"url"`
	} `json:"links"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetTripsTripIDLinksJSON400Response is a constructor method for a GetTripsTripIDLinks response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDLinksJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// PostTripsTripIDLinksJSON201Response is a constructor method for a PostTripsTripIDLinks response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDLinksJSON201Response(body struct {
	LinkID string `json:"linkId"`
}) *Response {
	return &Response{
		body:        body,
		Code:        201,
		contentType: "application/json",
	}
}

// PostTripsTripIDLinksJSON400Response is a constructor method for a PostTripsTripIDLinks response.
// A *Response is returned with the configured status code and content type from the spec.
func PostTripsTripIDLinksJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// GetTripsTripIDParticipantsJSON200Response is a constructor method for a GetTripsTripIDParticipants response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDParticipantsJSON200Response(body struct {
	Participants []struct {
		Email       openapi_types.Email `json:"email"`
		ID          string              `json:"id"`
		IsConfirmed bool                `json:"is_confirmed"`
		Name        *string             `json:"name"`
	} `json:"participants"`
}) *Response {
	return &Response{
		body:        body,
		Code:        200,
		contentType: "application/json",
	}
}

// GetTripsTripIDParticipantsJSON400Response is a constructor method for a GetTripsTripIDParticipants response.
// A *Response is returned with the configured status code and content type from the spec.
func GetTripsTripIDParticipantsJSON400Response(body Error) *Response {
	return &Response{
		body:        body,
		Code:        400,
		contentType: "application/json",
	}
}

// ServerInterface represents all server handlers.
type ServerInterface interface {
	// Confirms a participant on a trip.
	// (PATCH /participants/{participantId}/confirm)
	PatchParticipantsParticipantIDConfirm(w http.ResponseWriter, r *http.Request, participantID string) *Response
	// Create a new trip
	// (POST /trips)
	PostTrips(w http.ResponseWriter, r *http.Request) *Response
	// Get a trip details.
	// (GET /trips/{tripId})
	GetTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Update a trip.
	// (PUT /trips/{tripId})
	PutTripsTripID(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Get a trip activities.
	// (GET /trips/{tripId}/activities)
	GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Create a trip activity.
	// (POST /trips/{tripId}/activities)
	PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Confirm a trip and send e-mail invitations.
	// (GET /trips/{tripId}/confirm)
	GetTripsTripIDConfirm(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Invite someone to the trip.
	// (POST /trips/{tripId}/invites)
	PostTripsTripIDInvites(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Get a trip links.
	// (GET /trips/{tripId}/links)
	GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Create a trip link.
	// (POST /trips/{tripId}/links)
	PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request, tripID string) *Response
	// Get a trip participants.
	// (GET /trips/{tripId}/participants)
	GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request, tripID string) *Response
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler          ServerInterface
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

// PatchParticipantsParticipantIDConfirm operation middleware
func (siw *ServerInterfaceWrapper) PatchParticipantsParticipantIDConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "participantId" -------------
	var participantID string

	if err := runtime.BindStyledParameter("simple", false, "participantId", chi.URLParam(r, "participantId"), &participantID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "participantId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PatchParticipantsParticipantIDConfirm(w, r, participantID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostTrips operation middleware
func (siw *ServerInterfaceWrapper) PostTrips(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostTrips(w, r)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripID operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripID(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PutTripsTripID operation middleware
func (siw *ServerInterfaceWrapper) PutTripsTripID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PutTripsTripID(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripIDActivities operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripIDActivities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripIDActivities(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostTripsTripIDActivities operation middleware
func (siw *ServerInterfaceWrapper) PostTripsTripIDActivities(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostTripsTripIDActivities(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripIDConfirm operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripIDConfirm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripIDConfirm(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostTripsTripIDInvites operation middleware
func (siw *ServerInterfaceWrapper) PostTripsTripIDInvites(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostTripsTripIDInvites(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripIDLinks operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripIDLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripIDLinks(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// PostTripsTripIDLinks operation middleware
func (siw *ServerInterfaceWrapper) PostTripsTripIDLinks(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.PostTripsTripIDLinks(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

// GetTripsTripIDParticipants operation middleware
func (siw *ServerInterfaceWrapper) GetTripsTripIDParticipants(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// ------------- Path parameter "tripId" -------------
	var tripID string

	if err := runtime.BindStyledParameter("simple", false, "tripId", chi.URLParam(r, "tripId"), &tripID); err != nil {
		siw.ErrorHandlerFunc(w, r, &InvalidParamFormatError{err, "tripId"})
		return
	}

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := siw.Handler.GetTripsTripIDParticipants(w, r, tripID)
		if resp != nil {
			if resp.body != nil {
				render.Render(w, r, resp)
			} else {
				w.WriteHeader(resp.Code)
			}
		}
	})

	handler(w, r.WithContext(ctx))
}

type UnescapedCookieParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnescapedCookieParamError) Error() string {
	return fmt.Sprintf("error unescaping cookie parameter %s: %v", err.paramName, err.err)
}

func (err UnescapedCookieParamError) Unwrap() error { return err.err }

type UnmarshalingParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err UnmarshalingParamError) Error() string {
	return fmt.Sprintf("error unmarshaling parameter %s as JSON: %v", err.paramName, err.err)
}

func (err UnmarshalingParamError) Unwrap() error { return err.err }

type RequiredParamError struct {
	err       error
	paramName string
}

// Error implements error.
func (err RequiredParamError) Error() string {
	if err.err == nil {
		return fmt.Sprintf("query parameter %s is required, but not found", err.paramName)
	} else {
		return fmt.Sprintf("query parameter %s is required, but errored: %s", err.paramName, err.err)
	}
}

func (err RequiredParamError) Unwrap() error { return err.err }

type RequiredHeaderError struct {
	paramName string
}

// Error implements error.
func (err RequiredHeaderError) Error() string {
	return fmt.Sprintf("header parameter %s is required, but not found", err.paramName)
}

type InvalidParamFormatError struct {
	err       error
	paramName string
}

// Error implements error.
func (err InvalidParamFormatError) Error() string {
	return fmt.Sprintf("invalid format for parameter %s: %v", err.paramName, err.err)
}

func (err InvalidParamFormatError) Unwrap() error { return err.err }

type TooManyValuesForParamError struct {
	NumValues int
	paramName string
}

// Error implements error.
func (err TooManyValuesForParamError) Error() string {
	return fmt.Sprintf("expected one value for %s, got %d", err.paramName, err.NumValues)
}

// ParameterName is an interface that is implemented by error types that are
// relevant to a specific parameter.
type ParameterError interface {
	error
	// ParamName is the name of the parameter that the error is referring to.
	ParamName() string
}

func (err UnescapedCookieParamError) ParamName() string  { return err.paramName }
func (err UnmarshalingParamError) ParamName() string     { return err.paramName }
func (err RequiredParamError) ParamName() string         { return err.paramName }
func (err RequiredHeaderError) ParamName() string        { return err.paramName }
func (err InvalidParamFormatError) ParamName() string    { return err.paramName }
func (err TooManyValuesForParamError) ParamName() string { return err.paramName }

type ServerOptions struct {
	BaseURL          string
	BaseRouter       chi.Router
	ErrorHandlerFunc func(w http.ResponseWriter, r *http.Request, err error)
}

type ServerOption func(*ServerOptions)

// Handler creates http.Handler with routing matching OpenAPI spec.
func Handler(si ServerInterface, opts ...ServerOption) http.Handler {
	options := &ServerOptions{
		BaseURL:    "/",
		BaseRouter: chi.NewRouter(),
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	for _, f := range opts {
		f(options)
	}

	r := options.BaseRouter
	wrapper := ServerInterfaceWrapper{
		Handler:          si,
		ErrorHandlerFunc: options.ErrorHandlerFunc,
	}

	r.Route(options.BaseURL, func(r chi.Router) {
		r.Patch("/participants/{participantId}/confirm", wrapper.PatchParticipantsParticipantIDConfirm)
		r.Post("/trips", wrapper.PostTrips)
		r.Get("/trips/{tripId}", wrapper.GetTripsTripID)
		r.Put("/trips/{tripId}", wrapper.PutTripsTripID)
		r.Get("/trips/{tripId}/activities", wrapper.GetTripsTripIDActivities)
		r.Post("/trips/{tripId}/activities", wrapper.PostTripsTripIDActivities)
		r.Get("/trips/{tripId}/confirm", wrapper.GetTripsTripIDConfirm)
		r.Post("/trips/{tripId}/invites", wrapper.PostTripsTripIDInvites)
		r.Get("/trips/{tripId}/links", wrapper.GetTripsTripIDLinks)
		r.Post("/trips/{tripId}/links", wrapper.PostTripsTripIDLinks)
		r.Get("/trips/{tripId}/participants", wrapper.GetTripsTripIDParticipants)
	})
	return r
}

func WithRouter(r chi.Router) ServerOption {
	return func(s *ServerOptions) {
		s.BaseRouter = r
	}
}

func WithServerBaseURL(url string) ServerOption {
	return func(s *ServerOptions) {
		s.BaseURL = url
	}
}

func WithErrorHandler(handler func(w http.ResponseWriter, r *http.Request, err error)) ServerOption {
	return func(s *ServerOptions) {
		s.ErrorHandlerFunc = handler
	}
}

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"H4sIAAAAAAAC/+RZzW4bNxB+FYLtcWUprU8L9JA0QaAiKIwgRQ9BYIyXI4nJLrkluXIEQU/TQ0899gn8",
	"YsWQqxVXXssrJYot52JL1JDz/81wuOSZLkqtUDnL0yW32QwL8B+fZ07OpZPov4EQ0kmtIL8wukQT1ieQ",
	"W0x4GS0tObQ2SoeF//CjwQlP+Q/DDcdhzW5Y81rwVcLdokSecjAG/HcBDmn/RJsCHE/9wsDJAnlDbJ2R",
	"aspXq4Qb/KuSBgVP34etSSzQh2aLvvqImSMGDfP91JSiJVVVSXFboITrLKuMvSSafjok3EmXe5V3axf4",
	"edqYTZeKr4zR5l79BNrMyJJ+5yl/AYIRO7SOb+teoLUw7SHjmrBLqHdGlg8QZLTlVphtib0zYohYqokm",
	"fm2TvbIlZnIiM7j55+Y/tEwAe34xZiUYYJpdQfZpgErQMpR5IPtbszIHpc7QsEwr60x1868AJioDyiHT",
	"7Pc3f7LfdGUULmjnW519QmcR3Fnj/5Svz+AJn6OxQZ5nZ6OzkQ/CEhWUkqf8Z7+U8BLczFtsWIJxMpMl",
	"kK2W0bexWA0zrSbSFERYgstm9IHcAKTxWPCUX9DyRXRG9Hn88td6PzE0UKBDY3n6fsklyUdC8IQrKLwK",
	"MWseO8SZCpMam3ok3eoDbbalVjaEx0+jc/qXaeVQ+USE0tuftBh+tGSsZXQ+qqqgMFBVnlMA0H+4IjuT",
	"JD4A2o5/iROocsfe1lzJ5uej0V5Md4VuyN8OxnGS0q+2KgowC57y2vKWAYsMy7RiwJyRpQ8emJI3Ystb",
	"/oHOGRKJt12prevwurbunScJfkLrXmix2EvhPbJdoHVSQVB6yQup3qCauhlPz7e9n/DPg6ke4GdnYBAU",
	"XPI55DLUkSaqkkKqX8690bAAmdtLpy+lmstQbRpEaWLNU3XCdYwkvdkLOccknOllUGKPItGXiz9aXys0",
	"l4HV/Qr1VmAje2AQsnj5JZJaB8YdxwzbnUEUUDHfjSM6wqKladuu3VWijWCrW6j07Fj5Qtk77tOfbJml",
	"3tetzeOHPIPgkAFTeO0xLoK4gGcRtg2XQdkVCTLFDox7jQHi6M/4Za8SVtvv69au0TGj5HhAvDeoUVvV",
	"r6eW9rJuTFBEkHOldY6g+AFI0tlb9wGJligdeXM7v041u16jq1sHJtARNJ515FfCy6qrX6geLJdOvTk5",
	"XmPwmMrtYQX0+2vr/yhFqHHbPfzdBW7YvivXta7N8N1MWmZ05ZBdyzxnBl1lFIM8Z26GjHhadoXuGlH5",
	"FY8DjSMZKMFqVwbihOHck2pLR7qZrhzbCEKS76q20V39ROruLj9vTTpODu/bfltHXDwkIdzffVF8UL8e",
	"vQbsO+XbC6fvmAceDscbadeHP7L7Sx1ai0PuMNHeU7/HxMm3uDP1OhA/Gtv1uNrsM6Q7CtJ+t9O5xsdK",
	"MItUQwcFyJz5gYMXxfas8WFE0WdiF3w+rumfFAofbca1BTHfbubzFGI9hBqzukCtkDndNJB9ZtCbEM+l",
	"+mR7QtobT/v0RjaNDZr59Nd/sbzr+THhlWlnV2Vkv0nK+pWSDuh6Ctz5GBd0Pv2pidcjDvlasb6987eP",
	"6aMD9pd3tv2Csjfg02G34H539D5oo0whdEiTXO97Eg0y6dKVVh0lpFVk+lWS+Gn7CRaUbYscUlf6tl3N",
	"gP+Agf76dbHdAPUqPvVz3Vqmeyb29xSjdpty8jUpVmdXN7Za/R8AAP//rGRmniomAAA=",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
