package main

import (
	"github.com/fbuedding/iota-admin/internal/pkg/auth"
	fr "github.com/fbuedding/iota-admin/internal/pkg/fiwareRepository"
	"github.com/fbuedding/iota-admin/internal/pkg/sessionStore"
	"github.com/fbuedding/iota-admin/internal/server"
	"github.com/rs/zerolog/log"
)

func main() {
	repo, err := fr.NewFiwareRepo(fr.Sqlite)

	if err != nil {
		log.Fatal().Err(err).Msg("Could not creat sqlite repo")
	}
	s := server.New(auth.NewDebugAuth(), sessionStore.NewInMemory(), repo, 8080)
	s.Start()
}
