package gettitle

import (
	"bufio"
	"bytes"
	"strings"
	"testing"

	"github.com/elpinal/keepsake/log"
)

func TestGet(t *testing.T) {
	var buf bytes.Buffer
	logger := log.NewLogger(&buf, log.Debug)
	r := bufio.NewReader(strings.NewReader(`<title>The TITLE</title>`))

	s, err := Get(logger, r)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if expetced := `The TITLE`; s != expetced {
		t.Errorf("expected %q, but got %q", expetced, s)
		t.Logf("log: %s", buf.String())
	}
}

func TestGetEscape(t *testing.T) {
	var buf bytes.Buffer
	logger := log.NewLogger(&buf, log.Debug)
	r := bufio.NewReader(strings.NewReader("<title>Caf\u0026#xE9;</title>"))

	s, err := Get(logger, r)
	if err != nil {
		t.Errorf("error: %v", err)
	}

	if expetced := `Café`; s != expetced {
		t.Errorf("expected %q, but got %q", expetced, s)
		t.Logf("log: %s", buf.String())
	}
}
