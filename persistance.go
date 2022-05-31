package emailmanager

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// MailingEntryModel is a model of database table.
// MailingID and InsertTime are indexed to speed up queries.
// ContentHash has a unique index to prevent duplicated entries.
type MailingEntryModel struct {
	ID          uint      `gorm:"primarykey"`
	MailingID   MailingID `gorm:"index"`
	Email       string
	Title       string
	Content     string
	InsertTime  time.Time `gorm:"index"`
	ContentHash uint64    `gorm:"uniqueIndex"`
}

func ConnectDB() (*gorm.DB, error) {
	var err error
	var db *gorm.DB

	tries := 0
	maxTries := 3
	sleep := time.Second * 5

	// FIXME: dsn should be customizable
	dsn := "host=postgres user=postgres dbname=postgres password=postgres"

	for tries < maxTries {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		tries++
		time.Sleep(sleep)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	err = db.AutoMigrate(&MailingEntryModel{})
	if err != nil {
		return nil, fmt.Errorf("failed to auto migrate: %w", err)
	}

	return db, err
}
