package main

import (
	"github.com/fbuedding/iota-admin/internal/pkg/auth"
	"github.com/fbuedding/iota-admin/internal/pkg/sessionStore"
	"github.com/fbuedding/iota-admin/internal/server"
)

func main() {
  s := server.New(auth.NewDebugAuth(),sessionStore.NewInMemory(), 8080)
  s.Start() 
}
