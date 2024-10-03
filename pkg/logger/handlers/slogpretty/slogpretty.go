package slogpretty

import (
	"context"
	"encoding/json"
	"github.com/fatih/color"
	"io"
	stdLog "log"
	"log/slog"
	"os"
)

type PrettyHandlerOptions struct {
	SlogOpts *slog.HandlerOptions
}

type PrettyHandler struct {
	opts  PrettyHandlerOptions
	log   *stdLog.Logger
	attrs []slog.Attr
	slog.Handler
}

func (opts PrettyHandlerOptions) NewPrettyHandler(out io.Writer) *PrettyHandler {
	return &PrettyHandler{
		Handler: slog.NewJSONHandler(out, opts.SlogOpts),
		log:     stdLog.New(out, "", 0),
	}
}

// Handle processes a log record, applying formatting and color coding based on the log level.
// It constructs a formatted log message that includes the timestamp, log level, message, and any attributes.
//
// Params:
// - ctx (context.Context): The context of the log message (not used in this implementation).
// - rec (slog.Record): The log record to be processed and formatted.
//
// Returns:
// - An error if the log message could not be formatted or written.
func (h *PrettyHandler) Handle(_ context.Context, rec slog.Record) error {
	level := rec.Level.String() + ":"

	switch rec.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	fields := make(map[string]interface{}, rec.NumAttrs())

	rec.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()
		return true
	})

	for _, a := range h.attrs {
		fields[a.Key] = a.Value.Any()
	}

	var b []byte
	var err error

	if len(fields) > 0 {
		b, err = json.MarshalIndent(fields, "", "  ")
		if err != nil {
			return err
		}
	}

	timeStr := rec.Time.Format("[15:05:05.000]")
	msg := color.CyanString(rec.Message)

	h.log.Println(
		timeStr,
		level,
		msg,
		color.WhiteString(string(b)),
	)

	return nil
}

// WithAttrs returns a new instance of PrettyHandler with the specified additional attributes.
// This allows adding custom attributes to be included in every log message.
//
// Params:
// - attrs ([]slog.Attr): A slice of attributes to be included in each log message.
//
// Returns:
// - A new instance of PrettyHandler with the specified attributes.
func (h *PrettyHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &PrettyHandler{
		Handler: h.Handler,
		log:     h.log,
		attrs:   attrs,
	}
}

// SetupPrettySlog initializes a new slog.Logger with a pretty output handler.
// It configures the logger to output messages at the DEBUG level and formats them
// in a human-readable way with colors and indentation for easier reading.
//
// Returns:
// - *slog.Logger: A pointer to the configured slog.Logger instance.
func SetupPrettySlog() *slog.Logger {
	opts := PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
