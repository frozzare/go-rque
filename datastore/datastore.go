package datastore

import (
	"log"

	r "github.com/dancannon/gorethink"
	"github.com/frozzare/go-queue/shared"
)

var session *r.Session

// Instance will return Rethinkdb session instance
func Instance(c shared.Config) *r.Session {
	if session == nil {
		sen, err := r.Connect(r.ConnectOpts{
			Address:  c.Address,
			Database: c.Database,
		})

		if err != nil {
			log.Fatalf("Error connecting to DB: %s", err)
		}

		session = sen

		return session
	}

	return session
}
