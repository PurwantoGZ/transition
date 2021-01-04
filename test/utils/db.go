package utils

import (
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/fergusstrange/embedded-postgres"
	"github.com/phayes/freeport"

)

var port uint32

func init() {
	p,err := freeport.GetFreePort()
	if err != nil {
		panic(err)
	}
	port = uint32(p)
}

func StartEmbeddedPostgres() (close func() error,err error) {
	var postgres *embeddedpostgres.EmbeddedPostgres
	postgres = embeddedpostgres.NewDatabase(embeddedpostgres.
		DefaultConfig().
		Port(port).
		Username("postgres").
		Password("postgres").
		Database("order_svc"))
	err = postgres.Start()
	if err != nil {
		return nil,err
	}
	return postgres.Stop,nil
}

func TestDB() *pg.DB {
	return pg.Connect(&pg.Options{
		Addr: fmt.Sprintf("localhost:%d",port),
		Database: "order_svc",
		User: "postgres",
		Password: "postgres",
	})
}
