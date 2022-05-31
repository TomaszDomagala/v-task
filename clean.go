package emailmanager

import (
	"github.com/go-kit/kit/log"
	"gorm.io/gorm"
	"time"
)

var (
	cleanInterval = time.Minute * 5
)

func cleanOldMails(logger log.Logger, db *gorm.DB) error {
	tx := db.Where("insert_time < ?", time.Now().Add(-cleanInterval)).Delete(&MailingEntryModel{})
	if tx.Error != nil {
		logger.Log("method", "cleanOldMails", "error", tx.Error)
		return tx.Error
	}
	logger.Log("method", "cleanOldMails", "deleted", tx.RowsAffected)
	return nil
}

// StartCleanOldMails cleans old mails every 5 minutes
func StartCleanOldMails(logger log.Logger, db *gorm.DB) {
	for {
		if err := cleanOldMails(logger, db); err != nil {
			logger.Log("method", "StartCleanOldMails", "error", err)
		}
		time.Sleep(cleanInterval)
	}
}
