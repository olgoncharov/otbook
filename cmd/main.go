package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/olgoncharov/otbook/config"
	"github.com/olgoncharov/otbook/internal/app"
	"github.com/rs/zerolog/log"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewDefaultConfig()
	if err != nil {
		log.Fatal().Msgf("can't read config: %s", err.Error())
	}

	a, err := app.NewApp(cfg)
	if err != nil {
		log.Fatal().Msgf("can't init app: %s", err.Error())
	}
	a.Run(ctx)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig
	a.Shutdown(ctx)
}
