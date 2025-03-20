package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	dbFile := "youtube.sqlite"
	schemaDir := "schema"

	fmt.Printf("creating database: %s\n", dbFile)
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	fmt.Printf("getting schema from: %s\n", schemaDir)
	files, err := os.ReadDir(schemaDir)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to read directory: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("finished getting schema from: %s\n", schemaDir)

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			fmt.Printf("skipping non-sql file: %s\n", file.Name())
			continue
		}

		filePath := filepath.Join(schemaDir, file.Name())
		fmt.Printf("Executing %s\n", filePath)

		sqlBytes, err := os.ReadFile(filePath)
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to read SQL file %s: %s\n", filePath, err)
			os.Exit(1)
		}

		_, err = db.Exec(string(sqlBytes))
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "failed to execute SQL from %s: %s\n", filePath, err)
			os.Exit(1)
		}

		fmt.Printf("Finished executing %s\n", filePath)
	}
}
