package logger

import (
	"fmt"
	"io"
	"log/slog"
	"time"
)

func NewLogger(w io.Writer) *slog.Logger {
	th := slog.NewJSONHandler(
		w, &slog.HandlerOptions{
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				switch a.Key {
				case slog.SourceKey:
					src := a.Value.Any().(*slog.Source)
					return slog.String(slog.SourceKey, fmt.Sprintf("%s:%d", src.File, src.Line))
				case slog.TimeKey:
					t := a.Value.Time()
					return slog.String(slog.TimeKey, t.Format(time.Kitchen))
				}

				return a
			},
		},
	)

	return slog.New(th)
}
