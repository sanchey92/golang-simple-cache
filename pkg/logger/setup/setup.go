package logSetup

import (
	"github.com/sanchey92/golang-simple-cache/pkg/logger/handlers/slogpretty"
	"log/slog"
	"os"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// SetupLogger initializes a new slog.Logger based on the specified environment.
// It configures the logger to use different output formats and levels of verbosity
// depending on whether the application is running in local, development, or production mode.
//
// Params:
//   - env (string): The environment for which to set up the logger. Must be one of
//     the predefined environment constants (envLocal, envDev, envProd).
//
// Returns:
// - *slog.Logger: A pointer to the configured slog.Logger instance.
func SetupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slogpretty.SetupPrettySlog()
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}
