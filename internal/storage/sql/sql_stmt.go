package sql

import (
	"strconv"
	"unicode"

	"github.com/ivan-bokov/go-pdns/internal/stacktrace"
)

const (
	UNKNOWN = iota
	QUESTION
	DOLLAR
	NAMED
	AT
)

var defaultBinds = map[int][]string{
	DOLLAR:   []string{"postgres", "pgx", "pq-timeouts", "cloudsqlpostgres", "ql", "nrpostgres", "cockroach"},
	QUESTION: []string{"mysql", "sqlite3", "rqlite", "nrmysql", "nrsqlite3"},
	NAMED:    []string{"oci8", "ora", "goracle", "godror"},
	AT:       []string{"sqlserver"},
}

var binds map[string]int

func init() {
	binds = make(map[string]int)
	for bind, drivers := range defaultBinds {
		for _, driver := range drivers {
			binds[driver] = bind
		}
	}

}

func BindType(driverName string) int {
	itype, ok := binds[driverName]
	if !ok {
		return UNKNOWN
	}
	return itype
}

func CompileNamedQuery(queryString string, bindType int) (query string, names []string, err error) {
	qs := []rune(queryString)
	allowedBindRunes := []*unicode.RangeTable{unicode.Letter, unicode.Digit}
	names = make([]string, 0, 10)
	rebound := make([]rune, 0, len(qs))
	inName := false
	name := make([]rune, 0, 20)
	last := len(qs) - 1
	currentVar := 1

	for idx, b := range qs {
		switch {
		case b == ':':
			if inName && idx > 0 && qs[idx-1] == ':' {
				rebound = append(rebound, ':')
				inName = false
				continue
			} else if inName {
				err = stacktrace.New("unexpected `:` while reading named param at " + strconv.Itoa(idx))
				return
			}
			inName = true
			name = make([]rune, 0, 20)
		case inName && idx > 0 && b == '=' && len(name) == 0:
			rebound = append(rebound, ':', '=')
			inName = false
			continue
		case inName && (unicode.IsOneOf(allowedBindRunes, b) || b == '_' || b == '.') && idx != last:
			name = append(name, b)
		case inName:
			inName = false
			if idx == last && unicode.IsOneOf(allowedBindRunes, b) {
				name = append(name, b)
			}
			names = append(names, string(name))
			switch bindType {
			// oracle only supports named type bind vars even for positional
			case NAMED:
				rebound = append(rebound, ':')
				rebound = append(rebound, name...)
			case QUESTION, UNKNOWN:
				rebound = append(rebound, '?')
			case DOLLAR:
				rebound = append(rebound, '$')
				for _, b := range strconv.Itoa(currentVar) {
					rebound = append(rebound, b)
				}
				currentVar++
			case AT:
				rebound = append(rebound, '@', 'p')
				for _, b := range strconv.Itoa(currentVar) {
					rebound = append(rebound, b)
				}
				currentVar++
			}
			if idx != last {
				rebound = append(rebound, b)
			} else if !unicode.IsOneOf(allowedBindRunes, b) {
				rebound = append(rebound, b)
			}
		default:
			rebound = append(rebound, b)
		}
	}
	query = string(rebound)
	return
}
