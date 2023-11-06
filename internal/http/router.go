package http

import (
	"net/http"
	"strings"

	"github.com/ItsNotGoodName/radiomux/internal/openapi"
	"github.com/ItsNotGoodName/radiomux/pkg/echoext"
	"github.com/ItsNotGoodName/radiomux/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type EchoHandler interface {
	ServeEcho(c echo.Context) error
}

func NewRouter(
	androidWSServer EchoHandler,
	apiWSServer EchoHandler,
	apiServer openapi.ServerInterface,
	playerService http.Handler,
	presetService http.Handler,
	stateService http.Handler,
) *echo.Echo {
	e := echo.New()

	// Middleware
	e.Use(echoext.LoggerWithConfig(echoext.LoggerConfig{
		Format: []string{
			"id",
			"remote_ip",
			"host",
			"method",
			"user_agent",
			"status",
			"error",
			"latency_human",
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Skipper: func(c echo.Context) bool {
			// Prevent API 404's from being overwritten
			return strings.HasPrefix(c.Request().RequestURI, "/api")
		},
		Root:       "dist",
		Index:      "index.html",
		Browse:     false,
		HTML5:      true,
		Filesystem: web.DistFS(),
	}))

	// Routes
	e.GET("/ws", androidWSServer.ServeEcho) // TODO: remove this
	e.GET("/api/android/ws", androidWSServer.ServeEcho)
	e.GET("/api/ws", apiWSServer.ServeEcho)
	openapi.RegisterHandlers(e.Group("/api", middlewareOpenAPI), apiServer)
	e.Any("/rpc/PlayerService/*", echo.WrapHandler(playerService))
	e.Any("/rpc/PresetService/*", echo.WrapHandler(presetService))
	e.Any("/rpc/StateService/*", echo.WrapHandler(stateService))

	return e
}

func middlewareOpenAPI(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := next(c); err != nil {
			return openapi.ConvertErr(err)
		}
		return nil
	}
}
