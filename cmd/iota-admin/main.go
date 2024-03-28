package main

import (
	"os"

	"github.com/fbuedding/iota-admin/internal/globals"
	"github.com/fbuedding/iota-admin/internal/pkg/auth"
	bruteforceprotection "github.com/fbuedding/iota-admin/internal/pkg/bruteForceProtection"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/internal/pkg/sessionStore"
	"github.com/fbuedding/iota-admin/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	if globals.Conf.AppEnv == "development" {
		log.Logger = log.With().Caller().Logger()
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	log.Info().Msg("Startup...")
	if globals.Conf.BypassAuth {
		log.Warn().Msg("Authentication bypass is active")
	}

	repo, err := fr.NewFiwareRepo(fr.Sqlite)
	if err != nil {
		log.Fatal().Err(err).Msg("Could not creat sqlite repo")
	}
	s := server.New(getAuth(globals.Conf), sessionStore.NewInMemory(), bruteforceprotection.NewInMemory(globals.Conf.LoginAttempts), repo, 8080)
	s.Start()
}

func getAuth(conf globals.Config) auth.Authenticator {
	if conf.Username != "" && conf.Password != "" {
		return auth.NewUsernamePasswordAuth()
	} else if conf.AppEnv == "development" {
		return auth.NewDebugAuth()
	}
	log.Fatal().Msg("No authentication method defined or correctly defined")
	return nil
}
