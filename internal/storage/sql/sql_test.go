package sql

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindType(t *testing.T) {
	assert.Equal(t, BindType("sqlite3"), QUESTION)
	assert.NotEqual(t, BindType("postgres"), QUESTION)
	assert.Equal(t, BindType("rqlite"), QUESTION)
	assert.Equal(t, BindType("godror"), NAMED)
	assert.Equal(t, BindType("sqlserver"), AT)
	assert.Equal(t, BindType("cockroach"), DOLLAR)
}

func TestCompileNamedQuery(t *testing.T) {
	qr, names, err := CompileNamedQuery(`INSERT INTO foo (a,b,c,d) VALUES (:name, :age, :first, :last)`, QUESTION)
	actualNames := []string{"name", "age", "first", "last"}
	assert.Equal(t, err, nil)
	assert.Equal(t, qr, `INSERT INTO foo (a,b,c,d) VALUES (?, ?, ?, ?)`)
	assert.Equal(t, len(names), len(actualNames))
	for i, name := range names {
		assert.Equal(t, name, actualNames[i], fmt.Sprintf("expected %dth name to be %s, got %s", i+1, actualNames[i], name))
	}
}
