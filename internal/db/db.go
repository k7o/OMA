package db

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-yaml/yaml"
	"github.com/rs/zerolog/log"
)

func InitInMemoryDatabase(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	schema, err := schema()
	if err != nil {
		return nil, err
	}

	// Run migrations.
	if _, err := db.ExecContext(ctx, schema); err != nil {
		return nil, err
	}

	log.Debug().Msgf("Schema applied: \n %s", schema)

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

func schema() (string, error) {
	schema := ""

	migrationDirectories, err := sqlcMigrationDirs()
	if err != nil {
		return "", err
	}

	var files []string
	for _, path := range migrationDirectories {
		f, err := readMigrationsDir(path)
		if err != nil {
			return "", err
		}

		files = append(files, f...)
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

func sqlcMigrationDirs() ([]string, error) {
	sqlcFileName := "sqlc.yaml"
	sqlcFilePath, err := filepath.Abs(sqlcFileName)
	if err != nil {
		return nil, err
	}

	yamlFile, err := os.ReadFile(sqlcFilePath)
	if err != nil {
		return nil, err
	}

	// Parse yaml file
	var config struct {
		SQL []struct {
			Schema string `yaml:"schema"`
		} `yaml:"sql"`
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return nil, err
	}

	// Extract migration directories
	var migrationDirectories []string
	for _, item := range config.SQL {
		migrationDirectories = append(migrationDirectories, filepath.Join(strings.TrimSuffix(sqlcFilePath, sqlcFileName), item.Schema))
	}

	return migrationDirectories, nil
}
