package playgroundlogs

import (
	"embed"
)

//go:embed all:migrations/*.sql
var migrations embed.FS

// // Migrations returns the embedded filesystem for the playgroundlogs migrations.
// func (q *Queries) Migrations() (fs.FS, error) {
// 	return fs.Sub(migrations, "playgroundlogs/migrations")
// }

func (q *Queries) Migrations() embed.FS {
	return migrations
}
