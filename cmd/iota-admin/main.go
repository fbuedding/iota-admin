package main

import (
	"os"

	"github.com/fbuedding/iota-admin/internal/globals"
	"github.com/fbuedding/iota-admin/internal/pkg/auth"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/internal/pkg/sessionStore"
	"github.com/fbuedding/iota-admin/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.With().Caller().Logger()
	if globals.Conf.AppEnv == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if globals.Conf.BypassAuth {
		log.Warn().Msg("Authentication bypass is active")
	}

	repo, err := fr.NewFiwareRepo(fr.Sqlite)

	if err != nil {
		log.Fatal().Err(err).Msg("Could not creat sqlite repo")
	}
	s := server.New(auth.NewDebugAuth(), sessionStore.NewInMemory(), repo, 8080)
	s.Start()
}
