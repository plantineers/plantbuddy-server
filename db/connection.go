// Package db provides functions to connect to a database.
package db

import (
	"database/sql"

	// We know we'll never change the database, so it can be hard-coded here ...
	_ "github.com/mattn/go-sqlite3"
	"github.com/plantineers/plantbuddy-server/config"
)

// Session holds the database connection.
type Session struct {
	DB         *sql.DB
	Driver     string
	DataSource string
}

// NewSession creates a new session.
// The default database is configured in `buddy.json`. Required keys to use
// this function are `database.driverName` and `database.dataSource`.
// Note: as this project uses SQLite, the required driver is already imported.
func NewSession() *Session {
	conf := config.PlantBuddyConfig.Database

	return &Session{
		DB:         nil,
		Driver:     conf.DriverName,
		DataSource: conf.DataSource,
	}
}

// IsOpen returns true if the session is open.
func (s *Session) IsOpen() bool {
	return s.DB != nil
}

// Open opens a connection to the database.
func (s *Session) Open() error {
	db, err := sql.Open(s.Driver, s.DataSource)

	if err != nil {
		return err
	}

	s.DB = db
	return nil
}

// Close closes the connection to the database in case it is open.
func (s *Session) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
