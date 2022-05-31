package emailmanager

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"hash/fnv"
	"time"
)

type Service interface {
	AddMailingEntry(ctx context.Context, m MailingEntry) error
	// SendMails send mail to customers with given mailing_id
	SendMails(ctx context.Context, id MailingID) (int, error)
	// DeleteMails delete a message with given mailing_id
	DeleteMails(ctx context.Context, id MailingID) (int, error)
}

type MailingID int64

type MailingEntry struct {
	MailingID  MailingID
	Email      string
	Title      string
	Content    string
	InsertTime time.Time
}

func (entry MailingEntry) hash() (uint64, error) {
	h := fnv.New64()
	_, err := h.Write([]byte(fmt.Sprintf("%v", entry)))
	if err != nil {
		return 0, err
	}
	return h.Sum64(), nil
}

type persistentService struct {
	db *gorm.DB
	ms MailSender
}

func NewPersistentService(db *gorm.DB, ms MailSender) Service {
	return &persistentService{db, ms}
}

func (s *persistentService) AddMailingEntry(ctx context.Context, m MailingEntry) error {
	var err error
	var entry MailingEntryModel

	hash, err := m.hash()
	if err != nil {
		return fmt.Errorf("failed to hash mailing entry: %w", err)
	}

	entry.MailingID = m.MailingID
	entry.Email = m.Email
	entry.Title = m.Title
	entry.Content = m.Content
	entry.InsertTime = m.InsertTime
	entry.ContentHash = hash

	err = s.db.Create(&entry).Error

	if err != nil {
		return fmt.Errorf("failed to add mailing entry: %w", err)
	}

	return nil
}

func (s *persistentService) SendMails(ctx context.Context, id MailingID) (int, error) {
	var entries []MailingEntryModel
	// get mails not older than 5 minutes, as they should be treated as deleted
	t := time.Now().Add(-5 * time.Minute)
	// we delete entries before sending them, so we can't send them twice
	// however, if sending fails, it would be nice to somehow handle it
	// instead of deleting the entries, they could be flagged as sending
	// and on failure they could be flagged as failed

	tx := s.db.Clauses(clause.Returning{}).Where("mailing_id = ?", id).Where("insert_time > ?", t).Delete(&entries)

	if tx.Error != nil {
		return 0, fmt.Errorf("failed to get mailing entries: %w", tx.Error)
	}

	for _, entry := range entries {
		// could be done in parallel, but for now it is mocked anyway
		err := s.ms.SendMail(ctx, entry.Email, entry.Title, entry.Content)
		if err != nil {
			return 0, fmt.Errorf("failed to send mail: %w", err)
		}
	}

	return len(entries), nil
}

func (s *persistentService) DeleteMails(ctx context.Context, id MailingID) (int, error) {
	tx := s.db.Delete(MailingEntryModel{}, "mailing_id = ?", id)
	if tx.Error != nil {
		return 0, fmt.Errorf("failed to delete mailing entries: %w", tx.Error)
	}

	return int(tx.RowsAffected), nil
}
