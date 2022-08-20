package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/elpinal/keepsake/log"
	"github.com/elpinal/keepsake/server"
	"github.com/elpinal/keepsake/storage"
)

var port int = 7800

func main() {
	err := run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run() error {
	logger := log.NewLogger(os.Stdout, log.Debug)

	storage, err := storage.New(logger, "entries.db")
	if err != nil {
		return err
	}
	defer storage.Close()

	s := server.NewServer(logger, storage)
	http.Handle("/", s)
	http.Handle("/add", (*server.Add)(s))

	logger.LogInfo("Listening on port...", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)

	return err
}
