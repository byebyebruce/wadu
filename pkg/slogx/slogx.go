package slogx

import (
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
)

func init() {
	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stdout, &tint.Options{
			Level:     slog.LevelDebug,
			AddSource: true,

			TimeFormat: "2006/01/02 15:04:05.000",
		}),
	))
}
