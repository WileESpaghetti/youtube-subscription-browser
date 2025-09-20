package commands

import (
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jacob2161/sqlitebp"
)

type InitDBCmd struct {
	Schema string `help:"Path to the schema directory." type:"path" default:"db/migrations"`
}

func (ic *InitDBCmd) Run(ctx *Context) error {
	// TODO implement migrate command line options so we can repair/manage migrations

	db, err := sqlitebp.OpenReadWriteCreate(ctx.Database)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance("file://"+ic.Schema, "sqlite3", driver)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("initialized database: \"%s\"\n", ctx.Database)
	return nil
}
