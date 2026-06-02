package logger

import (
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewZapForwarder returns a *zap.Logger whose entries are forwarded into the
// panel's in-memory log buffer (the same one surfaced by the UI "Logs" view),
// mapping zap levels onto the panel logger. It lets libraries that log via zap
// (e.g. certmagic / ACME) show up in the panel logs. Entries below Info are
// dropped so certmagic's verbose debug output cannot flood the buffer.
func NewZapForwarder(name string) *zap.Logger {
	prefix := ""
	if name != "" {
		prefix = "[" + name + "] "
	}
	encCfg := zapcore.EncoderConfig{
		MessageKey:     "msg",
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	core := &forwardCore{
		LevelEnabler: zapcore.InfoLevel,
		enc:          zapcore.NewConsoleEncoder(encCfg),
		prefix:       prefix,
	}
	return zap.New(core)
}

// forwardCore is a zapcore.Core that renders each entry (message + fields) and
// writes it to the panel logger at the matching level.
type forwardCore struct {
	zapcore.LevelEnabler
	enc    zapcore.Encoder
	prefix string
}

func (c *forwardCore) With(fields []zapcore.Field) zapcore.Core {
	enc := c.enc.Clone()
	for i := range fields {
		fields[i].AddTo(enc)
	}
	return &forwardCore{
		LevelEnabler: c.LevelEnabler,
		enc:          enc,
		prefix:       c.prefix,
	}
}

func (c *forwardCore) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *forwardCore) Write(ent zapcore.Entry, fields []zapcore.Field) error {
	buf, err := c.enc.EncodeEntry(ent, fields)
	if err != nil {
		return err
	}
	msg := c.prefix + strings.TrimRight(buf.String(), "\n")
	buf.Free()

	switch {
	case ent.Level >= zapcore.ErrorLevel:
		Error(msg)
	case ent.Level == zapcore.WarnLevel:
		Warning(msg)
	default:
		Info(msg)
	}
	return nil
}

func (c *forwardCore) Sync() error { return nil }
