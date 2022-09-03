package log

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

type LogLevel int

const (
	Debug LogLevel = iota
	Info
	Warn
	Error
)

func (l LogLevel) String() string {
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

type LogLevelValue struct {
	Level LogLevel
}

func (p *LogLevelValue) String() string {
	if p != nil {
		return p.Level.String()
	}
	return ""
}

func (p *LogLevelValue) Set(s string) error {
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

func (l LogLevel) MarshalJSON() ([]byte, error) {
	bs := []byte{'"'}
	bs = append(bs, l.String()...)
	bs = append(bs, '"')
	return bs, nil
}

type Logger struct {
	level LogLevel
	out   io.Writer
}

func NewLogger(w io.Writer, level LogLevel) Logger {
	return Logger{
		level: level,
		out:   w,
	}
}

type Output struct {
	Date    time.Time
	Level   LogLevel
	Message string
	Value   any
}

func (l Logger) Log(level LogLevel, msg string, v any) {
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
