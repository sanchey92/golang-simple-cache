package main

import (
	"github.com/sanchey92/golang-simple-cache/internal/config"
	logSetup "github.com/sanchey92/golang-simple-cache/pkg/logger/setup"
)

func main() {
	cfg := config.MustLoad()
	log := logSetup.SetupLogger(cfg.Env)

	log.Info("Info log")
	log.Warn("Warn log")
	log.Error("Error log")
}
