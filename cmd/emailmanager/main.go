package main

import (
	"emailmanager"
	"fmt"
	"github.com/go-kit/kit/log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var logger log.Logger
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)
	logger = log.With(logger, "caller", log.DefaultCaller)

	db, err := emailmanager.ConnectDB()
	if err != nil {
		logger.Log("err", err)
		os.Exit(1)
	}
	logger.Log("msg", "Connected to DB")

	var ms emailmanager.MailSender
	ms = emailmanager.NewMockMailSender()

	var s emailmanager.Service
	s = emailmanager.NewPersistentService(db, ms)
	s = emailmanager.LoggingMiddleware(logger)(s)

	var h http.Handler
	h = emailmanager.MakeHTTPHandler(s, log.With(logger, "component", "HTTP"))

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		logger.Log("transport", "HTTP", "addr", ":8080")
		errs <- http.ListenAndServe(":8080", h)
	}()

	go emailmanager.StartCleanOldMails(logger, db)

	logger.Log("exit", <-errs)
}
