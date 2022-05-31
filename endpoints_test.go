package emailmanager

import (
	"context"
	"testing"
)

type mockService struct {
	_AddMailingEntry func(ctx context.Context, entry MailingEntry) error
	_SendMails       func(ctx context.Context, id MailingID) (int, error)
	_DeleteMails     func(ctx context.Context, id MailingID) (int, error)
}

func (s mockService) AddMailingEntry(ctx context.Context, entry MailingEntry) error {
	if s._AddMailingEntry != nil {
		return s._AddMailingEntry(ctx, entry)
	}
	return nil
}
func (s mockService) SendMails(ctx context.Context, id MailingID) (int, error) {
	if s._SendMails != nil {
		return s._SendMails(ctx, id)
	}
	return 0, nil
}
func (s mockService) DeleteMails(ctx context.Context, id MailingID) (int, error) {
	if s._DeleteMails != nil {
		return s._DeleteMails(ctx, id)
	}
	return 0, nil
}

func TestMakeSendMailsEndpoint_InvalidRequestType(t *testing.T) {
	var (
		ctx = context.Background()
		svc = mockService{}
		ep  = MakeSendMailsEndpoint(svc)
	)

	_, err := ep(ctx, "invalid")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

func TestMakeSendMailsEndpoint_ValidRequest(t *testing.T) {
	var svc mockService
	svc._SendMails = func(ctx context.Context, id MailingID) (int, error) {
		return 3, nil
	}
	ep := MakeSendMailsEndpoint(svc)

	res, err := ep(context.Background(), MailingID(1))
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	r, ok := res.(sendMailsResponse)
	if !ok {
		t.Errorf("expected sendMailsResponse, got %v", res)
	}
	if r.SendMails != 3 {
		t.Errorf("expected 3, got %d", r.SendMails)
	}
}
