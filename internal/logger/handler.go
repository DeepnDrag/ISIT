package logger

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/fatih/color"
	"io"
	"log/slog"
	"strings"
	"sync"
)

const timeFormat = "[2006-01-02 15:04:05]"

var (
	ErrWriteLog = errors.New("cannot write log")
)

type Handler struct {
	handler slog.Handler
	mu      *sync.Mutex
	out     io.Writer
}

type HandlerOpts struct {
	level slog.Level
	out   io.Writer
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return &Handler{handler: h.handler.WithAttrs(attrs), mu: h.mu}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{handler: h.handler.WithGroup(name), mu: h.mu}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = color.MagentaString(level)
	case slog.LevelInfo:
		level = color.BlueString(level)
	case slog.LevelWarn:
		level = color.YellowString(level)
	case slog.LevelError:
		level = color.RedString(level)
	}

	b, err := getAttrsJSON(r)
	if err != nil {
		return err
	}
	output := strings.Join([]string{color.WhiteString(r.Time.Format(timeFormat)), level,
		color.CyanString(r.Message), color.WhiteString(string(b)), "\n"}, " ")
	h.mu.Lock()
	defer h.mu.Unlock()
	_, err = h.out.Write([]byte(output))
	if err != nil {
		return errors.Join(ErrWriteLog, err)
	}

	return nil
}

func NewHandler(opts *HandlerOpts) *Handler {
	if opts == nil {
		opts = &HandlerOpts{}
	}
	var b bytes.Buffer
	return &Handler{
		handler: slog.NewJSONHandler(&b, &slog.HandlerOptions{
			Level: opts.level,
		}),
		mu:  &sync.Mutex{},
		out: opts.out,
	}
}

func getAttrsJSON(r slog.Record) ([]byte, error) {
	if r.NumAttrs() == 0 {
		return []byte{}, nil
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.Marshal(fields)
	if err != nil {
		return nil, err
	}

	return b, nil
}
