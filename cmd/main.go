package main

import (
	"github.com/sanchey92/golang-simple-cache/internal/config"
	"github.com/sanchey92/golang-simple-cache/internal/server"
	logSetup "github.com/sanchey92/golang-simple-cache/pkg/logger/setup"
	"log"
)

func main() {
	cfg := config.MustLoad()
	logger := logSetup.SetupLogger(cfg.Env)
	srv := server.New(cfg, logger)

	log.Fatal(srv.Start())
}
