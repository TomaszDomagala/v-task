package emailmanager

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Logging middleware for service, logs all requests

type Middleware func(Service) Service

func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) AddMailingEntry(ctx context.Context, m MailingEntry) error {
	defer func(begin time.Time) {
		mw.logger.Log("method", "AddMailingEntry", "mailing_id", m.MailingID, "email", m.Email, "title", m.Title, "insert_time", m.InsertTime, "took", time.Since(begin))
	}(time.Now())
	return mw.next.AddMailingEntry(ctx, m)
}

func (mw loggingMiddleware) SendMails(ctx context.Context, id MailingID) (int, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "SendMails", "mailing_id", id, "took", time.Since(begin))
	}(time.Now())
	return mw.next.SendMails(ctx, id)
}

func (mw loggingMiddleware) DeleteMails(ctx context.Context, id MailingID) (int, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeletedMails", "mailing_id", id, "took", time.Since(begin))
	}(time.Now())
	return mw.next.DeleteMails(ctx, id)
}
