// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.15.0 DO NOT EDIT.
package openapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/oapi-codegen/runtime"
)

// Defines values for EventType.
const (
	EventTypePlayerState        EventType = "player_state"
	EventTypePlayerStatePartial EventType = "player_state_partial"
)

// Defines values for PlayerPlaybackState.
const (
	Buffering PlayerPlaybackState = 2
	Ended     PlayerPlaybackState = 4
	Idle      PlayerPlaybackState = 1
	Ready     PlayerPlaybackState = 3
)

// Error defines model for Error.
type Error struct {
	Message string `json:"message"`
}

// Event defines model for Event.
type Event struct {
	union json.RawMessage
}

// EventBase defines model for EventBase.
type EventBase struct {
	Type EventType `json:"type"`
}

// EventDataPlayerState defines model for EventDataPlayerState.
type EventDataPlayerState struct {
	Data []PlayerState `json:"data"`
	Type EventType     `json:"type"`
}

// EventDataPlayerStatePartial defines model for EventDataPlayerStatePartial.
type EventDataPlayerStatePartial struct {
	Data []PlayerStatePartial `json:"data"`
	Type EventType            `json:"type"`
}

// EventType defines model for EventType.
type EventType string

// Player defines model for Player.
type Player struct {
	Id    int64  `json:"id"`
	Name  string `json:"name"`
	Token string `json:"token"`
}

// PlayerPlaybackState defines model for PlayerPlaybackState.
type PlayerPlaybackState int

// PlayerState defines model for PlayerState.
type PlayerState struct {
	Connected     bool                `json:"connected"`
	Genre         string              `json:"genre"`
	Id            int64               `json:"id"`
	Loading       bool                `json:"loading"`
	MaxVolume     int                 `json:"max_volume"`
	MinVolume     int                 `json:"min_volume"`
	Muted         bool                `json:"muted"`
	Name          string              `json:"name"`
	PlaybackError string              `json:"playback_error"`
	PlaybackState PlayerPlaybackState `json:"playback_state"`
	Playing       bool                `json:"playing"`
	Ready         bool                `json:"ready"`
	Station       string              `json:"station"`
	Title         string              `json:"title"`
	Uri           string              `json:"uri"`
	Volume        int                 `json:"volume"`
}

// PlayerStatePartial defines model for PlayerStatePartial.
type PlayerStatePartial struct {
	Connected     *bool                `json:"connected,omitempty"`
	Genre         *string              `json:"genre,omitempty"`
	Id            int64                `json:"id"`
	Loading       *bool                `json:"loading,omitempty"`
	MaxVolume     *int                 `json:"max_volume,omitempty"`
	MinVolume     *int                 `json:"min_volume,omitempty"`
	Muted         *bool                `json:"muted,omitempty"`
	Name          *string              `json:"name,omitempty"`
	PlaybackError *string              `json:"playback_error,omitempty"`
	PlaybackState *PlayerPlaybackState `json:"playback_state,omitempty"`
	Playing       *bool                `json:"playing,omitempty"`
	Ready         *bool                `json:"ready,omitempty"`
	Station       *string              `json:"station,omitempty"`
	Title         *string              `json:"title,omitempty"`
	Uri           *string              `json:"uri,omitempty"`
	Volume        *int                 `json:"volume,omitempty"`
}

// Preset defines model for Preset.
type Preset struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

// PostPlayersJSONBody defines parameters for PostPlayers.
type PostPlayersJSONBody struct {
	Name string `json:"name"`
}

// PostPlayersIdJSONBody defines parameters for PostPlayersId.
type PostPlayersIdJSONBody struct {
	Name *string `json:"name,omitempty"`
}

// PostPlayersIdMediaParams defines parameters for PostPlayersIdMedia.
type PostPlayersIdMediaParams struct {
	Uri string `form:"uri" json:"uri"`
}

// PostPlayersIdPresetParams defines parameters for PostPlayersIdPreset.
type PostPlayersIdPresetParams struct {
	Preset int64 `form:"preset" json:"preset"`
}

// PostPlayersIdVolumeParams defines parameters for PostPlayersIdVolume.
type PostPlayersIdVolumeParams struct {
	Delta *int  `form:"delta,omitempty" json:"delta,omitempty"`
	Mute  *bool `form:"mute,omitempty" json:"mute,omitempty"`
}

// PostPresetsJSONBody defines parameters for PostPresets.
type PostPresetsJSONBody struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

// PostPresetsIdJSONBody defines parameters for PostPresetsId.
type PostPresetsIdJSONBody struct {
	Name *string `json:"name,omitempty"`
	Url  *string `json:"url,omitempty"`
}

// PostPlayersJSONRequestBody defines body for PostPlayers for application/json ContentType.
type PostPlayersJSONRequestBody PostPlayersJSONBody

// PostPlayersIdJSONRequestBody defines body for PostPlayersId for application/json ContentType.
type PostPlayersIdJSONRequestBody PostPlayersIdJSONBody

// PostPresetsJSONRequestBody defines body for PostPresets for application/json ContentType.
type PostPresetsJSONRequestBody PostPresetsJSONBody

// PostPresetsIdJSONRequestBody defines body for PostPresetsId for application/json ContentType.
type PostPresetsIdJSONRequestBody PostPresetsIdJSONBody

// AsEventDataPlayerState returns the union data inside the Event as a EventDataPlayerState
func (t Event) AsEventDataPlayerState() (EventDataPlayerState, error) {
	var body EventDataPlayerState
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromEventDataPlayerState overwrites any union data inside the Event as the provided EventDataPlayerState
func (t *Event) FromEventDataPlayerState(v EventDataPlayerState) error {
	v.Type = "player_state"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeEventDataPlayerState performs a merge with any union data inside the Event, using the provided EventDataPlayerState
func (t *Event) MergeEventDataPlayerState(v EventDataPlayerState) error {
	v.Type = "player_state"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

// AsEventDataPlayerStatePartial returns the union data inside the Event as a EventDataPlayerStatePartial
func (t Event) AsEventDataPlayerStatePartial() (EventDataPlayerStatePartial, error) {
	var body EventDataPlayerStatePartial
	err := json.Unmarshal(t.union, &body)
	return body, err
}

// FromEventDataPlayerStatePartial overwrites any union data inside the Event as the provided EventDataPlayerStatePartial
func (t *Event) FromEventDataPlayerStatePartial(v EventDataPlayerStatePartial) error {
	v.Type = "player_state_partial"
	b, err := json.Marshal(v)
	t.union = b
	return err
}

// MergeEventDataPlayerStatePartial performs a merge with any union data inside the Event, using the provided EventDataPlayerStatePartial
func (t *Event) MergeEventDataPlayerStatePartial(v EventDataPlayerStatePartial) error {
	v.Type = "player_state_partial"
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}

	merged, err := runtime.JsonMerge(t.union, b)
	t.union = merged
	return err
}

func (t Event) Discriminator() (string, error) {
	var discriminator struct {
		Discriminator string `json:"type"`
	}
	err := json.Unmarshal(t.union, &discriminator)
	return discriminator.Discriminator, err
}

func (t Event) ValueByDiscriminator() (interface{}, error) {
	discriminator, err := t.Discriminator()
	if err != nil {
		return nil, err
	}
	switch discriminator {
	case "player_state":
		return t.AsEventDataPlayerState()
	case "player_state_partial":
		return t.AsEventDataPlayerStatePartial()
	default:
		return nil, errors.New("unknown discriminator value: " + discriminator)
	}
}

func (t Event) MarshalJSON() ([]byte, error) {
	b, err := t.union.MarshalJSON()
	return b, err
}

func (t *Event) UnmarshalJSON(b []byte) error {
	err := t.union.UnmarshalJSON(b)
	return err
}

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (DELETE /players)
	DeletePlayers(ctx echo.Context) error

	// (GET /players)
	GetPlayers(ctx echo.Context) error

	// (POST /players)
	PostPlayers(ctx echo.Context) error

	// (DELETE /players/{id})
	DeletePlayersId(ctx echo.Context, id int64) error

	// (GET /players/{id})
	GetPlayersId(ctx echo.Context, id int64) error

	// (POST /players/{id})
	PostPlayersId(ctx echo.Context, id int64) error

	// (POST /players/{id}/media)
	PostPlayersIdMedia(ctx echo.Context, id int64, params PostPlayersIdMediaParams) error

	// (POST /players/{id}/pause)
	PostPlayersIdPause(ctx echo.Context, id int64) error

	// (POST /players/{id}/play)
	PostPlayersIdPlay(ctx echo.Context, id int64) error

	// (POST /players/{id}/preset)
	PostPlayersIdPreset(ctx echo.Context, id int64, params PostPlayersIdPresetParams) error

	// (POST /players/{id}/seek)
	PostPlayersIdSeek(ctx echo.Context, id int64) error

	// (POST /players/{id}/stop)
	PostPlayersIdStop(ctx echo.Context, id int64) error

	// (POST /players/{id}/volume)
	PostPlayersIdVolume(ctx echo.Context, id int64, params PostPlayersIdVolumeParams) error

	// (DELETE /presets)
	DeletePresets(ctx echo.Context) error

	// (GET /presets)
	GetPresets(ctx echo.Context) error

	// (POST /presets)
	PostPresets(ctx echo.Context) error

	// (DELETE /presets/{id})
	DeletePresetsId(ctx echo.Context, id int64) error

	// (GET /presets/{id})
	GetPresetsId(ctx echo.Context, id int64) error

	// (POST /presets/{id})
	PostPresetsId(ctx echo.Context, id int64) error
}

// ServerInterfaceWrapper converts echo contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler ServerInterface
}

// DeletePlayers converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePlayers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeletePlayers(ctx)
	return err
}

// GetPlayers converts echo context to params.
func (w *ServerInterfaceWrapper) GetPlayers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPlayers(ctx)
	return err
}

// PostPlayers converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayers(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayers(ctx)
	return err
}

// DeletePlayersId converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePlayersId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeletePlayersId(ctx, id)
	return err
}

// GetPlayersId converts echo context to params.
func (w *ServerInterfaceWrapper) GetPlayersId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPlayersId(ctx, id)
	return err
}

// PostPlayersId converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayersId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayersId(ctx, id)
	return err
}

// PostPlayersIdMedia converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayersIdMedia(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PostPlayersIdMediaParams
	// ------------- Required query parameter "uri" -------------

	err = runtime.BindQueryParameter("form", true, true, "uri", ctx.QueryParams(), &params.Uri)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter uri: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayersIdMedia(ctx, id, params)
	return err
}

// PostPlayersIdPause converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayersIdPause(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayersIdPause(ctx, id)
	return err
}

// PostPlayersIdPlay converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayersIdPlay(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayersIdPlay(ctx, id)
	return err
}

// PostPlayersIdPreset converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayersIdPreset(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PostPlayersIdPresetParams
	// ------------- Required query parameter "preset" -------------

	err = runtime.BindQueryParameter("form", true, true, "preset", ctx.QueryParams(), &params.Preset)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter preset: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayersIdPreset(ctx, id, params)
	return err
}

// PostPlayersIdSeek converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayersIdSeek(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayersIdSeek(ctx, id)
	return err
}

// PostPlayersIdStop converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayersIdStop(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayersIdStop(ctx, id)
	return err
}

// PostPlayersIdVolume converts echo context to params.
func (w *ServerInterfaceWrapper) PostPlayersIdVolume(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Parameter object where we will unmarshal all parameters from the context
	var params PostPlayersIdVolumeParams
	// ------------- Optional query parameter "delta" -------------

	err = runtime.BindQueryParameter("form", true, false, "delta", ctx.QueryParams(), &params.Delta)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter delta: %s", err))
	}

	// ------------- Optional query parameter "mute" -------------

	err = runtime.BindQueryParameter("form", true, false, "mute", ctx.QueryParams(), &params.Mute)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter mute: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPlayersIdVolume(ctx, id, params)
	return err
}

// DeletePresets converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePresets(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeletePresets(ctx)
	return err
}

// GetPresets converts echo context to params.
func (w *ServerInterfaceWrapper) GetPresets(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPresets(ctx)
	return err
}

// PostPresets converts echo context to params.
func (w *ServerInterfaceWrapper) PostPresets(ctx echo.Context) error {
	var err error

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPresets(ctx)
	return err
}

// DeletePresetsId converts echo context to params.
func (w *ServerInterfaceWrapper) DeletePresetsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.DeletePresetsId(ctx, id)
	return err
}

// GetPresetsId converts echo context to params.
func (w *ServerInterfaceWrapper) GetPresetsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.GetPresetsId(ctx, id)
	return err
}

// PostPresetsId converts echo context to params.
func (w *ServerInterfaceWrapper) PostPresetsId(ctx echo.Context) error {
	var err error
	// ------------- Path parameter "id" -------------
	var id int64

	err = runtime.BindStyledParameterWithLocation("simple", false, "id", runtime.ParamLocationPath, ctx.Param("id"), &id)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid format for parameter id: %s", err))
	}

	// Invoke the callback with all the unmarshaled arguments
	err = w.Handler.PostPresetsId(ctx, id)
	return err
}

// This is a simple interface which specifies echo.Route addition functions which
// are present on both echo.Echo and echo.Group, since we want to allow using
// either of them for path registration
type EchoRouter interface {
	CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
}

// RegisterHandlers adds each server route to the EchoRouter.
func RegisterHandlers(router EchoRouter, si ServerInterface) {
	RegisterHandlersWithBaseURL(router, si, "")
}

// Registers handlers, and prepends BaseURL to the paths, so that the paths
// can be served under a prefix.
func RegisterHandlersWithBaseURL(router EchoRouter, si ServerInterface, baseURL string) {

	wrapper := ServerInterfaceWrapper{
		Handler: si,
	}

	router.DELETE(baseURL+"/players", wrapper.DeletePlayers)
	router.GET(baseURL+"/players", wrapper.GetPlayers)
	router.POST(baseURL+"/players", wrapper.PostPlayers)
	router.DELETE(baseURL+"/players/:id", wrapper.DeletePlayersId)
	router.GET(baseURL+"/players/:id", wrapper.GetPlayersId)
	router.POST(baseURL+"/players/:id", wrapper.PostPlayersId)
	router.POST(baseURL+"/players/:id/media", wrapper.PostPlayersIdMedia)
	router.POST(baseURL+"/players/:id/pause", wrapper.PostPlayersIdPause)
	router.POST(baseURL+"/players/:id/play", wrapper.PostPlayersIdPlay)
	router.POST(baseURL+"/players/:id/preset", wrapper.PostPlayersIdPreset)
	router.POST(baseURL+"/players/:id/seek", wrapper.PostPlayersIdSeek)
	router.POST(baseURL+"/players/:id/stop", wrapper.PostPlayersIdStop)
	router.POST(baseURL+"/players/:id/volume", wrapper.PostPlayersIdVolume)
	router.DELETE(baseURL+"/presets", wrapper.DeletePresets)
	router.GET(baseURL+"/presets", wrapper.GetPresets)
	router.POST(baseURL+"/presets", wrapper.PostPresets)
	router.DELETE(baseURL+"/presets/:id", wrapper.DeletePresetsId)
	router.GET(baseURL+"/presets/:id", wrapper.GetPresetsId)
	router.POST(baseURL+"/presets/:id", wrapper.PostPresetsId)

}
