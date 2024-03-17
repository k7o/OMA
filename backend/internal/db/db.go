package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"sort"
)

func InitInMemoryDatabase(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	schema, err := Schema("./internal/decisionlogs/migrations")
	if err != nil {
		return nil, err
	}

	// Run migrations.
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, err
	}

	log.Println("Schema applied: \n", schema)

	return db, nil
}

func InitDatabase() (*sql.DB, error) {
	if _, err := os.Stat("./data/1234/data.db"); os.IsNotExist(err) {
		os.MkdirAll("./data/1234", 0755)
		os.Create("./data/1234/data.db")
	}

	db, err := sql.Open("sqlite3", "./data/1234/data.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}

func readMigrationsDir(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	var migrationFiles []string
	for _, file := range files {
		if !file.IsDir() {
			migrationFiles = append(migrationFiles, filepath.Join(path, file.Name()))
		}
	}

	sort.Strings(migrationFiles)
	return migrationFiles, nil
}

func Schema(path string) (string, error) {
	schema := ""
	files, err := readMigrationsDir(path)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			return "", err
		}
		schema += string(content)
	}

	return schema, nil
}
