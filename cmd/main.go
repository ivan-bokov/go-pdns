package main

import (
	"os"

	"github.com/ivan-bokov/go-pdns/internal/handler"
	"github.com/ivan-bokov/go-pdns/internal/service"
	"github.com/ivan-bokov/go-pdns/internal/storage/sqlite"
)

func main() {
	err := os.Remove("sql.db")
	if err != nil {
		panic(err)
	}
	storage := sqlite.New("sql.db")
	err = storage.CreateTable()
	if err != nil {
		panic(err)
	}
	svc := service.New(storage, true)
	handlerHTTP := handler.New(svc)
	err = handlerHTTP.InitRoutes().Run()
	if err != nil {
		panic(err)
	}
}
