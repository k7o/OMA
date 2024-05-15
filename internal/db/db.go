package db

import (
	"context"
	"database/sql"
	"io/fs"
	"oma/contract"
	"os"
	"path/filepath"
	"sort"

	"github.com/rs/zerolog/log"
)

func InitInMemoryDatabase(ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}
	return db, nil
}

func Migrate(ctx context.Context, db *sql.DB, migrationEmbeds ...contract.MigrationEmbed) error {
	schema, err := schema(migrationEmbeds)
	if err != nil {
		return err
	}

	_, err = db.Exec(schema)
	if err != nil {
		return err
	}

	log.Info().Msg("migrated database")
	log.Debug().Msgf("Database schema: \n%s", schema)

	return nil
}

func schema(migrationEmbeds []contract.MigrationEmbed) (string, error) {
	schemaFiles := map[string]string{}
	for _, migration := range migrationEmbeds {
		fsys := migration.Migrations()
		err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// Check if the current entry is not a directory
			if !d.IsDir() {
				fileBytes, err := fs.ReadFile(fsys, path)
				if err != nil {
					return err
				}
				schemaFiles[filepath.Base(path)] = string(fileBytes)
			}
			return nil
		})
		if err != nil {
			return "", err
		}
	}

	fileNames := make([]string, 0, len(schemaFiles))
	for fileName := range schemaFiles {
		fileNames = append(fileNames, fileName)
	}

	sort.Strings(fileNames)

	var schema string
	for _, fileName := range fileNames {
		schema += schemaFiles[fileName] + "\n"
	}

	return schema, nil
}

func InitDatabase() (*sql.DB, error) {
	if _, err := os.Stat("./data/1234/data.db"); os.IsNotExist(err) {
		os.MkdirAll("./data/1234", 0755)
		os.Create("./data/1234/data.db")
	}

	db, err := sql.Open("sqlite", "./data/1234/data.db")
	if err != nil {
		return nil, err
	}

	return db, nil
}
