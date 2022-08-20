package sentry

import (
	"runtime"
	"strings"
	"time"

	raven "github.com/getsentry/sentry-go"
	"go.uber.org/zap/zapcore"
)

const (
	_platform          = "go"
	_traceContextLines = 3
	_traceSkipFrames   = 2
)

func ravenSeverity(lvl zapcore.Level) raven.Level {
	switch lvl {
	case zapcore.DebugLevel:
		return raven.LevelDebug
	case zapcore.InfoLevel:
		return raven.LevelInfo
	case zapcore.WarnLevel:
		return raven.LevelWarning
	case zapcore.ErrorLevel:
		return raven.LevelError
	case zapcore.DPanicLevel:
		return raven.LevelFatal
	case zapcore.PanicLevel:
		return raven.LevelFatal
	case zapcore.FatalLevel:
		return raven.LevelFatal
	default:
		// Unrecognized levels are fatal.
		return raven.LevelFatal
	}
}

type trace struct {
	Disabled bool
}

// Configuration is a minimal set of parameters for Sentry integration.
type Configuration struct {
	DSN   string `yaml:"DSN"`
	Trace trace
}

// Build uses the provided configuration to construct a Sentry-backed logging core.
func (c Configuration) Build() (zapcore.Core, error) {
	//client, err := raven.New(c.DSN)
	err := raven.Init(raven.ClientOptions{
		Dsn:              c.DSN,
		AttachStacktrace: !c.Trace.Disabled,
		BeforeSend: func(event *raven.Event, hint *raven.EventHint) *raven.Event {
			// Modify the event here
			if errMsg, ok := event.Extra["error"].(string); ok && strings.Contains(errMsg, "OAuthException") {
				event.Extra["error"] = strings.Replace(errMsg, "OAuthException", "OAfbException", -1)
			}
			return event
		},
	})

	if err != nil {
		return zapcore.NewNopCore(), err
	}
	return newCore(c, zapcore.ErrorLevel), nil
}

type core struct {
	zapcore.LevelEnabler
	trace

	fields map[string]interface{}
}

func newCore(cfg Configuration, enab zapcore.LevelEnabler) *core {
	sentryCore := &core{
		LevelEnabler: enab,
		trace:        cfg.Trace,
		fields:       make(map[string]interface{}),
	}
	return sentryCore
}

func (c *core) With(fs []zapcore.Field) zapcore.Core {
	return c.with(fs)
}

func (c *core) Check(ent zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if c.Enabled(ent.Level) {
		return ce.AddCore(ent, c)
	}
	return ce
}

func (c *core) Write(ent zapcore.Entry, fs []zapcore.Field) error {
	clone := c.with(fs)

	packet := raven.NewEvent()
	packet.Message = ent.Message
	packet.Timestamp = ent.Time
	packet.Level = ravenSeverity(ent.Level)
	packet.Platform = _platform
	packet.Extra = clone.fields
	packet.Extra["runtime.Version"] = runtime.Version()
	packet.Extra["runtime.NumCPU"] = runtime.NumCPU()
	//packet := &raven.Packet{
	//	Message:   ent.Message,
	//	Timestamp: raven.Timestamp(ent.Time),
	//	Level:     ravenSeverity(ent.Level),
	//	Platform:  _platform,
	//	Extra:     clone.fields,
	//}

	//if !c.trace.Disabled {
	//	trace := raven.NewStacktrace()
	//	if trace != nil {
	//		packet.Interfaces = append(packet.Interfaces, trace)
	//	}
	//}
	raven.CaptureEvent(packet)
	// We may be crashing the program, so should flush any buffered events.
	if ent.Level > zapcore.ErrorLevel {
		raven.Flush(2 * time.Second)
	}
	return nil
}

func (c *core) Sync() error {
	raven.Flush(2 * time.Second)
	return nil
}

func (c *core) with(fs []zapcore.Field) *core {
	// Copy our map.
	m := make(map[string]interface{}, len(c.fields))
	for k, v := range c.fields {
		m[k] = v
	}

	// Add fields to an in-memory encoder.
	enc := zapcore.NewMapObjectEncoder()
	for _, f := range fs {
		f.AddTo(enc)
	}

	// Merge the two maps.
	for k, v := range enc.Fields {
		m[k] = v
	}

	return &core{
		LevelEnabler: c.LevelEnabler,
		trace:        c.trace,
		fields:       m,
	}
}
