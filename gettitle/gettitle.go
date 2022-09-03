package gettitle

import (
	"bufio"
	"html"
	"io"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/elpinal/keepsake/log"
)

func Get(logger log.Logger, r *bufio.Reader) (string, error) {
	p := NewParser(logger, r)
	s, err := p.Parse()
	if err == io.EOF {
		return "", nil
	}
	return s, err
}

func FromURL(logger log.Logger, url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	for _, ty := range resp.Header["Content-Type"] {
		logger.LogDebug("content type", ty)
		if ty == "application/pdf" {
			logger.LogDebug("PDF detected", nil)
			return "", nil
		}
	}

	r := bufio.NewReader(resp.Body)
	return Get(logger, r)
}

type Parser struct {
	logger log.Logger
	*bufio.Reader
}

func NewParser(logger log.Logger, r *bufio.Reader) *Parser {
	return &Parser{
		logger: logger,
		Reader: r,
	}
}

func (p *Parser) Parse() (string, error) {
	// TODO: Support timeout.
	// TODO: Trim whitespace.
	for {
		err := p.skipTo('<')
		if err != nil {
			return "", err
		}

		p.logPeek(1)

		err = p.skipWhitespace()
		if err != nil {
			return "", err
		}

		p.logPeek(1)

		found, err := p.exactLower("title")
		if err != nil {
			return "", err
		}
		if !found {
			p.logger.LogDebug("'title' not found, continuing", nil)
			continue
		}
		p.logger.LogDebug("'title' found", nil)

		p.logPeek(1)

		err = p.skipTo('>')
		if err != nil {
			return "", err
		}

		p.logPeek(1)

		s, err := p.readTo('<')
		if err != nil {
			return "", err
		}
		return s, nil
	}
}

func (p *Parser) logPeek(n int) {
	bs, err := p.Peek(n)
	if err != nil {
		p.logger.LogDebug("peek error", err.Error())
	} else {
		p.logger.LogDebug("peek", string(bs))
	}
}

func (p *Parser) skipTo(to byte) error {
	for {
		b, err := p.ReadByte()
		if err != nil {
			return err
		}
		if b == to {
			return nil
		}
	}
}

func (p *Parser) skipWhitespace() error {
	for {
		b, err := p.ReadByte()
		if err != nil {
			return err
		}
		if b != ' ' {
			return p.UnreadByte()
		}
	}
}

func (p *Parser) exactLower(s string) (bool, error) {
	var notFound bool

	bs := make([]byte, len(s))
	_, err := p.Read(bs)
	if err != nil {
		return notFound, err
	}

	s1 := strings.ToLower(string(bs))
	if s == s1 {
		return true, nil
	}
	return notFound, nil
}

func (p *Parser) readTo(b byte) (string, error) {
	bs, err := p.ReadBytes(b)
	if err != nil {
		return "", err
	}
	bs = bs[:len(bs)-1]
	s := string(bs)
	if !utf8.ValidString(s) {
		p.logger.LogWarn("not utf8 string", bs)
	}
	return html.UnescapeString(s), nil
}
