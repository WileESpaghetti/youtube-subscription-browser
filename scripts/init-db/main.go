package main

import (
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jacob2161/sqlitebp"
)

const defaultDatabaseFile = "youtube-migrations.sqlite"
const defaultMigrationsDir = "db/migrations"

func main() {
	// TODO convert database and migrations directory to absolute paths
	var dbFile string
	if len(os.Args) > 1 {
		dbFile = os.Args[1]
	} else {
		dbFile = defaultDatabaseFile
	}

	schemaDir := defaultMigrationsDir

	db, err := sqlitebp.OpenReadWriteCreate(dbFile)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
	if err != nil {
		log.Fatal(err)
	}
	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance("file://"+schemaDir, "sqlite3", driver)
	if err != nil {
		log.Fatal(err)
	}
	err = m.Up()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Database initialized.")
}
