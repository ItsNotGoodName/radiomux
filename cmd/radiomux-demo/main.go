package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/ItsNotGoodName/radiomux/internal/android"
	"github.com/ItsNotGoodName/radiomux/internal/androidmock"
	"github.com/ItsNotGoodName/radiomux/internal/api"
	"github.com/ItsNotGoodName/radiomux/internal/apiws"
	"github.com/ItsNotGoodName/radiomux/internal/build"
	"github.com/ItsNotGoodName/radiomux/internal/bus"
	"github.com/ItsNotGoodName/radiomux/internal/demo"
	"github.com/ItsNotGoodName/radiomux/internal/http"
	"github.com/ItsNotGoodName/radiomux/internal/rpc"
	"github.com/ItsNotGoodName/radiomux/internal/webrpc"
	"github.com/ItsNotGoodName/radiomux/pkg/sutureext"
	"github.com/Rican7/lieut"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/thejerf/suture/v4"
)

func main() {
	ctx := context.Background()

	cfg := demo.New()

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

func run(cfg *demo.Config) lieut.Executor {
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
		bus := bus.New()

		// Store
		playerStore := demo.NewPlayerStore()
		presetStore := demo.NewPresetStore()

		// Services
		androidStatePubSub := android.NewStateMemPubSub()
		androidStateStore := android.NewStateMemStore(androidStatePubSub, bus, playerStore)
		androidStateService := android.NewStateService(androidStatePubSub, androidStateStore)
		androidController := android.NewController(androidStateService, bus)
		androidWSServer := demo.NewAndroidWSServer()
		apiWSServer := apiws.NewServer(androidStateService, playerStore)
		apiServer := api.NewServer(playerStore, androidWSServer)
		playerService := webrpc.
			NewPlayerServiceServer(rpc.
				NewPlayerService(playerStore, androidWSServer))
		presetService := webrpc.
			NewPresetServiceServer(rpc.
				NewPresetService(presetStore))
		stateService := webrpc.
			NewStateServiceServer(
				demo.NewStateService(
					rpc.NewStateService(androidController, presetStore)))

		for _, mockPlayer := range demo.MockPlayers {
			mock := androidmock.NewMock(mockPlayer.ID, androidController, androidStateService)
			super.Add(mock)
		}

		// HTTP
		httpRouter := http.NewRouter(
			androidWSServer,
			apiWSServer,
			apiServer,
			playerService,
			presetService,
			stateService,
		)
		httpServer := http.NewServer(httpRouter, fmt.Sprintf("%s:%d", cfg.Host, cfg.Port))
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
