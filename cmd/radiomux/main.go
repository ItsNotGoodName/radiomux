package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/androidws"
	"github.com/ItsNotGoodName/radiomux/internal/api"
	"github.com/ItsNotGoodName/radiomux/internal/apiws"
	"github.com/ItsNotGoodName/radiomux/internal/build"
	"github.com/ItsNotGoodName/radiomux/internal/bus"
	"github.com/ItsNotGoodName/radiomux/internal/config"
	"github.com/ItsNotGoodName/radiomux/internal/file"
	"github.com/ItsNotGoodName/radiomux/internal/http"
	"github.com/ItsNotGoodName/radiomux/internal/rpc"
	"github.com/ItsNotGoodName/radiomux/internal/webrpc"
	"github.com/ItsNotGoodName/radiomux/pkg/echoext"
	"github.com/ItsNotGoodName/radiomux/pkg/sutureext"
	"github.com/ItsNotGoodName/radiomux/web"
	"github.com/Rican7/lieut"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func main() {
	ctx := context.Background()

	cfg := config.New()

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cfg.WithFlag(flags)

	app := lieut.NewSingleCommandApp(
		lieut.AppInfo{
			Name:    "radiomux",
			Version: build.Current.Version,
		},
		run(cfg),
		flags,
		os.Stdout,
		os.Stderr,
	)

	code := app.Run(ctx, os.Args[1:])

	os.Exit(code)
}

func run(cfg *config.Config) lieut.Executor {
	return func(ctx context.Context, arguments []string) error {
		err := cfg.Parse()
		if err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		// Supervisor
		super := suture.New("root", suture.Spec{
			EventHook: sutureext.EventHook(),
		})

		// Bus
		bus, err := bus.New()
		if err != nil {
			return fmt.Errorf("failed to create bus: %w", err)
		}

		// Store
		jsonStore := file.NewStore(cfg.File)
		playerStore := file.NewPlayerStore(jsonStore, bus)
		presetStore := file.NewPresetStore(jsonStore, bus)

		// Services
		androidStatePubSub := android.NewStateMemPubSub()
		androidStateStore, close1 := android.NewStateMemStore(androidStatePubSub, bus)
		defer close1()
		androidStateService := android.NewStateService(androidStatePubSub, androidStateStore)
		androidController, close2 := android.NewController(androidStateService, bus)
		defer close2()
		androidWSServer := androidws.NewServer(playerStore, androidController, androidStateService, cfg.HTTPURL)
		apiWSServer := apiws.NewServer(androidStateService, playerStore)
		apiServer := api.NewServer(playerStore, androidWSServer)
		playerService := webrpc.
			NewPlayerServiceServer(rpc.
				NewPlayerService(playerStore, androidWSServer))
		presetService := webrpc.
			NewPresetServiceServer(rpc.
				NewPresetService(presetStore))
		stateService := webrpc.
			NewStateServiceServer(rpc.
				NewStateService(androidController, presetStore))

		// DEBUG START
		// for i := 1; i <= 10; i++ {
		// 	mock := androidmock.NewMock(int64(i), androidController, androidStateService)
		// 	super.Add(mock)
		// }
		// DEBUG END

		// HTTP
		e := echo.New()

		// - Middleware
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
				return strings.HasPrefix(c.Request().RequestURI, "/api")
			},
			Root:       "dist",
			Index:      "index.html",
			Browse:     false,
			HTML5:      true,
			Filesystem: web.DistFS(),
		}))

		// - Routes
		e.GET("/ws", androidWSServer.Handle)
		e.GET(androidws.Path, androidWSServer.Handle)
		e.GET("/api/ws", apiWSServer.Handle)
		api.MountServer(e, apiServer)
		e.Any("/rpc/PlayerService/*", echo.WrapHandler(playerService))
		e.Any("/rpc/PresetService/*", echo.WrapHandler(presetService))
		e.Any("/rpc/StateService/*", echo.WrapHandler(stateService))

		httpServer := http.NewServer(e, fmt.Sprintf("%s:%d", cfg.HTTPHost, cfg.HTTPPort))
		super.Add(httpServer)

		return super.Serve(ctx)
	}
}

var (
	builtBy    = "unknown"
	commit     = ""
	date       = ""
	version    = "dev"
	repoURL    = "https://github.com/ItsNotGoodName/radiomux"
	releaseURL = ""
)

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	build.Current = build.Build{
		BuiltBy:    builtBy,
		Commit:     commit,
		Date:       date,
		Version:    version,
		RepoURL:    repoURL,
		ReleaseURL: releaseURL,
	}
}
