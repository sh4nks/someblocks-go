package gormstore

import (
	"errors"
	"log"
	"time"

	"gorm.io/gorm"
)

// Based on: https://github.com/alexedwards/scs/pull/84

// GORMStore represents the session store.
type GORMStore struct {
	db          *gorm.DB
	stopCleanup chan bool
}

type session struct {
	Token  string `gorm:"primary_key;type:varchar(100)"`
	Data   []byte
	Expiry time.Time `gorm:"index"`
}

// New returns a new GORMStore instance, with a background cleanup goroutine
// that runs every 5 minutes to remove expired session data.
func New(db *gorm.DB) (*GORMStore, error) {
	return NewWithCleanupInterval(db, 5*time.Minute)
}

// NewWithCleanupInterval returns a new GORMStore instance. The cleanupInterval
// parameter controls how frequently expired session data is removed by the
// background cleanup goroutine. Setting it to 0 prevents the cleanup goroutine
// from running (i.e. expired sessions will not be removed).
func NewWithCleanupInterval(db *gorm.DB, cleanupInterval time.Duration) (*GORMStore, error) {
	p := &GORMStore{db: db}
	if err := p.migrate(); err != nil {
		return nil, err
	}
	if cleanupInterval > 0 {
		go p.startCleanup(cleanupInterval)
	}
	return p, nil
}

// Find returns the data for a given session token from the GORMStore instance.
// If the session token is not found or is expired, the returned exists flag will
// be set to false.
func (p *GORMStore) Find(token string) ([]byte, bool, error) {
	row := &session{}
	result := p.db.First(row, "token = ? AND expiry >= ?", token, time.Now())
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, false, nil
	}
	return row.Data, true, nil
}

// Commit adds a session token and data to the GORMStore instance with the
// given expiry time. If the session token already exists, then the data and expiry
// time are updated.
func (p *GORMStore) Commit(token string, b []byte, expiry time.Time) error {
	row := &session{}
	result := p.db.Where(session{Token: token}).Assign(session{Data: b, Expiry: expiry}).FirstOrCreate(&row)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// Delete removes a session token and corresponding data from the GORMStore
// instance.
func (p *GORMStore) Delete(token string) error {
	result := p.db.Delete(&session{}, "token = ?", token)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (p *GORMStore) migrate() error {
	var tableOptions string
	// Set table options for MySQL database dialect
	if p.db.Dialector.Name() == "mysql" {
		tableOptions = "ENGINE=InnoDB CHARSET=utf8mb4"
	}

	p.db.Set("gorm:table_options", tableOptions).AutoMigrate(&session{})
	return nil
}

func (p *GORMStore) startCleanup(interval time.Duration) {
	p.stopCleanup = make(chan bool)
	ticker := time.NewTicker(interval)
	for {
		select {
		case <-ticker.C:
			err := p.deleteExpired()
			if err != nil {
				log.Println(err)
			}
		case <-p.stopCleanup:
			ticker.Stop()
			return
		}
	}
}

// StopCleanup terminates the background cleanup goroutine for the GORMStore
// instance. It's rare to terminate this; generally GORMStore instances and
// their cleanup goroutines are intended to be long-lived and run for the lifetime
// of your application.
//
// There may be occasions though when your use of the GORMStore is transient.
// An example is creating a new GORMStore instance in a test function. In this
// scenario, the cleanup goroutine (which will run forever) will prevent the
// GORMStore object from being garbage collected even after the test function
// has finished. You can prevent this by manually calling StopCleanup.
func (p *GORMStore) StopCleanup() {
	if p.stopCleanup != nil {
		p.stopCleanup <- true
	}
}

func (p *GORMStore) deleteExpired() error {
	result := p.db.Delete(&session{}, "expiry < ?", time.Now())
	if result.Error != nil {
		return result.Error
	}
	return nil
}
