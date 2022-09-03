package log

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type Level int

const (
	Debug Level = iota
	Info
	Warn
	Error
)

func (l Level) String() string {
	switch l {
	case Debug:
		return "DEBUG"
	case Info:
		return "INFO"
	case Warn:
		return "WARN"
	case Error:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

type LevelValue struct {
	Level Level
}

func (p *LevelValue) String() string {
	if p != nil {
		return p.Level.String()
	}
	return ""
}

func (p *LevelValue) Set(s string) error {
	switch strings.ToLower(s) {
	case "debug":
		p.Level = Debug
	case "info":
		p.Level = Info
	case "warn":
		p.Level = Warn
	case "error":
		p.Level = Error
	default:
		return fmt.Errorf("unknown log level: %q", s)
	}

	return nil
}

func (l Level) MarshalJSON() ([]byte, error) {
	bs := []byte{'"'}
	bs = append(bs, l.String()...)
	bs = append(bs, '"')
	return bs, nil
}

type Logger struct {
	level Level
	out   io.Writer
}

func NewLogger(w io.Writer, level Level) Logger {
	return Logger{
		level: level,
		out:   w,
	}
}

type Output struct {
	Date    time.Time
	Level   Level
	Message string
	Value   any
}

func (l Logger) Log(level Level, msg string, v any) {
	if l.level > level {
		return
	}
	t := time.Now()
	e := json.NewEncoder(l.out)
	err := e.Encode(Output{
		Date:    t,
		Level:   level,
		Message: msg,
		Value:   v,
	})
	if err != nil {
		fmt.Fprintln(l.out, msg)
		fmt.Fprintf(l.out, "Logger.Log: failed to encode to json: %v\n", err)
	}
}

func (l Logger) LogDebug(msg string, v any) {
	l.Log(Debug, msg, v)
}

func (l Logger) LogInfo(msg string, v any) {
	l.Log(Info, msg, v)
}

func (l Logger) LogWarn(msg string, v any) {
	l.Log(Warn, msg, v)
}

func (l Logger) LogError(msg string, v any) {
	l.Log(Error, msg, v)
}
