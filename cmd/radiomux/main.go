package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/androidws"
	"github.com/ItsNotGoodName/radiomux/internal/api"
	"github.com/ItsNotGoodName/radiomux/internal/apiws"
	"github.com/ItsNotGoodName/radiomux/internal/build"
	"github.com/ItsNotGoodName/radiomux/internal/bus"
	"github.com/ItsNotGoodName/radiomux/internal/file"
	"github.com/ItsNotGoodName/radiomux/internal/http"
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

	flags := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	filePath := flags.String("file", "radiomux.json", "File path to JSON database.")

	app := lieut.NewSingleCommandApp(
		lieut.AppInfo{
			Name:    "radiomux",
			Version: build.Current.Version,
		},
		run(filePath),
		flags,
		os.Stdout,
		os.Stderr,
	)

	code := app.Run(ctx, os.Args[1:])

	os.Exit(code)
}

func run(filePath *string) lieut.Executor {
	return func(ctx context.Context, arguments []string) error {
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
		jsonStore := file.NewStore(*filePath)
		playerStore := file.NewPlayerStore(jsonStore, bus)
		presetStore := file.NewPresetStore(jsonStore, bus)

		// Services
		androidStatePubSub := android.NewStateMemPubSub()
		androidStateStore, c1 := android.NewStateMemStore(androidStatePubSub, bus)
		defer c1()
		androidStateService := android.NewStateService(androidStatePubSub, androidStateStore)
		androidController, c2 := android.NewController(androidStateService, bus)
		defer c2()
		androidWSServer := androidws.NewServer(playerStore, androidController, androidStateService)
		apiWSServer := apiws.NewServer(androidStateService, playerStore)
		apiServer := api.NewServer(playerStore, presetStore, androidController)

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

		// - Routes
		e.GET("/ws", androidWSServer.Handle)
		e.GET("/api/ws", apiWSServer.Handle)
		api.MountServer(e, apiServer, "/api")
		err = echoext.MountFS(e, web.FS())
		if err != nil {
			return fmt.Errorf("failed to mount web filesystem: %w", err)
		}

		httpServer := http.NewServer(e, ":8080")
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
