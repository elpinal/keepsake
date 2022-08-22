package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/elpinal/keepsake/entry"
	"github.com/elpinal/keepsake/gettitle"
	"github.com/elpinal/keepsake/log"
)

type Storage interface {
	Add(url string, title string, date time.Time) error
	List(limit int, offset int) ([]entry.Entry, error)
	Count() (int, error)
	Export(enc *json.Encoder, limit int, offset int) error
	Import(dec *json.Decoder) error
}

type Server struct {
	logger  log.Logger
	storage Storage
	dev     bool
}

func NewServer(logger log.Logger, storage Storage) *Server {
	return &Server{
		logger:  logger,
		storage: storage,
		dev:     true, // TODO: inherit from CLI options.
	}
}

func updateTemplate() *template.Template {
	return template.Must(template.New("index.html").ParseFiles("./resources/index.html"))
}

var tmpl = updateTemplate()

func (s *Server) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	s.logger.LogInfo("request path", req.URL.Path)

	count, err := s.storage.Count()
	s.logger.LogInfo("the number of entries", count)

	entries, err := s.storage.List(100, 0) // TODO
	if err != nil {
		s.logger.LogError("storage.List", err.Error())
	}

	if s.dev {
		tmpl = updateTemplate()
	}

	err = tmpl.Execute(w, entries)
	if err != nil {
		s.logger.LogError("tmpl.Execute", err.Error())
	}
	end := time.Now()
	s.logger.LogInfo("duration", end.Sub(start).String())
}

type Add Server

func (s *Add) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	switch req.Method {
	case http.MethodPost:
		url := req.PostFormValue("url")
		s.logger.LogInfo("POST", map[string]string{"url": url})

		title, err := gettitle.FromURL(s.logger, url)
		if err != nil {
			s.logger.LogError("gettitle.FromURL", err.Error())
			break
		}
		s.logger.LogInfo("title", title)
		if title == "" {
			title = url
		}

		err = s.storage.Add(url, title, start)
		if err != nil {
			s.logger.LogError("storage.Add", err.Error())
		}
	}

	s.logger.LogInfo("/add: redirecting to /", nil)
	http.Redirect(w, req, "/", http.StatusSeeOther)

	end := time.Now()
	s.logger.LogInfo("duration", end.Sub(start).String())
}

type Export Server

func (s *Export) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	s.logger.LogInfo("/export", nil)

	enc := json.NewEncoder(w)
	err := s.storage.Export(enc, 1000, 0) // TODO
	if err != nil {
		s.logger.LogError("storage.Export", err.Error())
	}

	end := time.Now()
	s.logger.LogInfo("/export: duration", end.Sub(start).String())
}

type Import Server

func (s *Import) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	start := time.Now()
	s.logger.LogInfo("/import", nil)

	if req.Method != http.MethodPost {
		s.logger.LogInfo("/import: not POST", req.Method)
		http.Redirect(w, req, "/", http.StatusSeeOther)
		return
	}

	f, header, err := req.FormFile("input")
	if err != nil {
		s.logger.LogError("FormFile", err.Error())
	}
	s.logger.LogInfo("/import: header", header)

	dec := json.NewDecoder(f)
	err = s.storage.Import(dec)
	if err != nil {
		s.logger.LogError("storage.Import", err.Error())
	}

	s.logger.LogInfo("/import: redirecting to /", nil)
	http.Redirect(w, req, "/", http.StatusSeeOther)

	end := time.Now()
	s.logger.LogInfo("/import: duration", end.Sub(start).String())
}
