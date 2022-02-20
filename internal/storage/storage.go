package storage

type IStorage interface {
	CreateTable() error
	Query(stmt string, args ...interface{}) (IResult, error)
	Exec(stmt string, args ...interface{}) (int, error)
	Close()
}

type IResult interface {
	Next() bool
	Err() error
	Columns() ([]string, error)
	Scan(dest ...interface{}) error
}
