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
	"github.com/ItsNotGoodName/radiomux/internal/config"
	"github.com/ItsNotGoodName/radiomux/internal/core"
	"github.com/ItsNotGoodName/radiomux/internal/http"
	"github.com/ItsNotGoodName/radiomux/internal/jsondb"
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

		core.Init(cfg.HTTPURL)

		// Supervisor
		super := suture.New("root", suture.Spec{
			EventHook: sutureext.EventHook(),
		})

		// Bus
		bus := bus.New()
		androidStatePubSub := android.NewStateMemPubSub()

		// Store
		jsonStore := jsondb.NewStore(cfg.File)
		playerStore := jsondb.NewPlayerStore(jsonStore, bus)
		presetStore := jsondb.NewPresetStore(jsonStore, bus)
		sourceStore := jsondb.NewSourceStore(jsonStore)
		androidStateStore := android.NewStateMemStore(androidStatePubSub, bus, playerStore)

		// Services
		androidStateService := android.NewStateService(androidStatePubSub, androidStateStore)
		androidController := android.NewController(androidStateService, bus)
		androidWSServer := androidws.NewServer(playerStore, androidController, androidStateService, cfg.HTTPURLRaw)
		notificationServer := apiws.NewNotificationServer(bus, playerStore)
		apiWSServer := apiws.NewServer(androidStateService, playerStore, notificationServer)
		apiServer := api.NewServer(playerStore, sourceStore)
		playerService := webrpc.
			NewPlayerServiceServer(rpc.
				NewPlayerService(playerStore))
		presetService := webrpc.
			NewPresetServiceServer(rpc.
				NewPresetService(presetStore))
		stateService := webrpc.
			NewStateServiceServer(rpc.
				NewStateService(androidController, presetStore))

		// Bootstrap
		if _, err = androidStateStore.Sync(ctx); err != nil {
			return fmt.Errorf("failed to sync android state store: %w", err)
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
		httpServer := http.NewServer(httpRouter, fmt.Sprintf("%s:%d", cfg.HTTPHost, cfg.HTTPPort))
		super.Add(httpServer)

		return super.Serve(ctx)
	}
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}
