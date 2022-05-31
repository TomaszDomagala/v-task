package emailmanager

import (
	"context"
	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	AddMailingEntryEndpoint endpoint.Endpoint
	SendMailsEndpoint       endpoint.Endpoint
	DeleteMailsEndpoint     endpoint.Endpoint
}

func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		AddMailingEntryEndpoint: MakeAddMailingEntryEndpoint(s),
		SendMailsEndpoint:       MakeSendMailsEndpoint(s),
		DeleteMailsEndpoint:     MakeDeleteMailsEndpoint(s),
	}
}

func MakeAddMailingEntryEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		entry, ok := request.(MailingEntry)
		if !ok {
			return nil, ErrInvalidArgumentType
		}
		err := s.AddMailingEntry(ctx, entry)
		return addMailingEntryResponse{}, err
	}
}

func MakeSendMailsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, ok := request.(MailingID)
		if !ok {
			return nil, ErrInvalidArgumentType
		}
		n, err := s.SendMails(ctx, id)
		return sendMailsResponse{n}, err
	}
}

func MakeDeleteMailsEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		id, ok := request.(MailingID)
		if !ok {
			return nil, ErrInvalidArgumentType
		}
		n, err := s.DeleteMails(ctx, id)
		return deleteMailsResponse{n}, err
	}
}

type addMailingEntryResponse struct {
}

type sendMailsResponse struct {
	SendMails int `json:"send_mails"`
}

type deleteMailsResponse struct {
	DeletedMails int `json:"deleted_mails"`
}

type addMailingEntryRequest struct {
	MailingID  MailingID `json:"mailing_id"`
	Email      string    `json:"email"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	InsertTime string    `json:"insert_time"`
}

type sendMailsRequest struct {
	MailingID MailingID `json:"mailing_id"`
}
