package decisionlogs

import (
	"embed"
)

//go:embed all:migrations/*.sql
var migrations embed.FS

// // Migrations returns the embedded filesystem for the decisionlogs migrations.
// func (q *Queries) Migrations() (fs.FS, error) {
// 	return fs.Sub(migrations, "playgroundlogs/migrations")
// }

// Migrations returns the embedded filesystem for the decisionlogs migrations.
func (q *Queries) Migrations() embed.FS {
	return migrations
}
