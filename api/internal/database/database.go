package database

import (
	"context"
	"errors"
	"log"
	"net"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// RetryDelay is how long to wait after a timeout before the single retry.
// Allows time for the database to become reachable after a timeout.
const RetryDelay = 30 * time.Second

// Connect opens a GORM connection. If the first attempt times out, it waits
// RetryDelay and tries once more.
func Connect(dsn string) (*gorm.DB, error) {
	return connect(dsn, RetryDelay)
}

func connect(dsn string, retryDelay time.Duration) (*gorm.DB, error) {
	open := func() (*gorm.DB, error) {
		return gorm.Open(postgres.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logger.Warn),
		})
	}

	db, err := open()
	if err != nil && isTimeout(err) {
		log.Printf("database connect timed out; retrying in %s", retryDelay)
		time.Sleep(retryDelay)
		db, err = open()
	}
	return db, err
}

func isTimeout(err error) bool {
	var ne net.Error
	return errors.Is(err, context.DeadlineExceeded) || (errors.As(err, &ne) && ne.Timeout())
}
