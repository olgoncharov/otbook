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

	a := app.NewApp(cfg)
	a.Run()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-sig
	a.Shutdown(ctx)
}
