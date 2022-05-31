package emailmanager

import (
	"context"
	"encoding/json"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/transport"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"

	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	router := mux.NewRouter()
	endpoints := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorHandler(transport.NewLogErrorHandler(logger)),
		httptransport.ServerErrorEncoder(encodeError),
	}

	router.Methods("POST").Path("/api/messages").Handler(httptransport.NewServer(
		endpoints.AddMailingEntryEndpoint,
		decodeAddMailingEntryRequest,
		encodeResponse,
		options...,
	))
	router.Methods("POST").Path("/api/messages/send").Handler(httptransport.NewServer(
		endpoints.SendMailsEndpoint,
		decodeSendMailsRequest,
		encodeResponse,
		options...,
	))
	router.Methods("DELETE").Path("/api/messages/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteMailsEndpoint,
		decodeDeleteMailsRequest,
		encodeResponse,
		options...,
	))

	return router
}

func encodeResponse(ctx context.Context, writer http.ResponseWriter, response interface{}) error {
	writer.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(writer).Encode(response)
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	switch err {
	default:
		w.WriteHeader(http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func addMailingEntryRequestToMailingEntry(req addMailingEntryRequest) (MailingEntry, error) {
	var entry MailingEntry
	var err error

	entry.MailingID = req.MailingID
	entry.Email = req.Email
	entry.Title = req.Title
	entry.Content = req.Content

	entry.InsertTime, err = time.Parse(time.RFC3339Nano, req.InsertTime)

	return entry, err
}

func decodeAddMailingEntryRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	var entryReq addMailingEntryRequest
	if err = json.NewDecoder(req.Body).Decode(&entryReq); err != nil {
		return nil, err
	}
	return addMailingEntryRequestToMailingEntry(entryReq)
}

func decodeSendMailsRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	var sendReq sendMailsRequest
	if err = json.NewDecoder(req.Body).Decode(&sendReq); err != nil {
		return nil, err
	}
	return sendReq.MailingID, nil
}

func decodeDeleteMailsRequest(_ context.Context, req *http.Request) (request interface{}, err error) {
	vars := mux.Vars(req)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	mailingID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, ErrMailingIDNaN
	}

	return MailingID(mailingID), nil
}
