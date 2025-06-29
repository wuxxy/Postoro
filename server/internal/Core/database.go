package Core

import (
	"errors"
	"fmt"
	"server/internal/User"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database struct {
	mu   sync.RWMutex
	conn *gorm.DB
	dsn  string
}
type DatabaseInfo struct {
	Host string `json:"host"`
	Port string `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Name string `json:"name"`
	SSL  bool   `json:"ssl"`
}

// NewDatabase creates a new, empty instance (not connected yet)
func NewDatabase() *Database {
	return &Database{}
}

// Connect sets or resets the DB connection using the provided DSN
func (d *Database) Connect(dsn string) error {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("failed to connect to DB: %w", err)
	}
	// Auto-migrate your models
	err = db.AutoMigrate(&User.User{}, &User.Session{}) // you can add more: &Post{}, &Thread{}
	if err != nil {
		panic("failed to migrate database: " + err.Error())
	}
	d.mu.Lock()
	d.conn = db
	d.dsn = dsn
	d.mu.Unlock()

	return nil
}

// Reconnect attempts to reconnect using the last known DSN
func (d *Database) Reconnect() error {
	d.mu.RLock()
	dsn := d.dsn
	d.mu.RUnlock()

	if dsn == "" {
		return errors.New("no previous DSN stored")
	}

	return d.Connect(dsn)
}

// Get returns the underlying *gorm.DB or nil
func (d *Database) Get() *gorm.DB {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.conn
}

// IsConnected returns true if DB is connected
func (d *Database) IsConnected() bool {
	return d.Get() != nil
}

var DB = NewDatabase()
