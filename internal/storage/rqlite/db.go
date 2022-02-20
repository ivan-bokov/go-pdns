package sqlite

import "github.com/ivan-bokov/go-pdns/internal/stacktrace"

type Rqlite struct {
}

func New() *Rqlite {
	return &Rqlite{}
}

func (db *Rqlite) Get(stmt string, args ...interface{}) error {
	return stacktrace.New("No implement")
}

func (db *Rqlite) Set(stmt string, args ...interface{}) error {
	return stacktrace.New("No implement")
}
