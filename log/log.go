package log

import (
	"encoding/json"
	"fmt"
	"io"
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
