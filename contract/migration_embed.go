package contract

import "embed"

type MigrationEmbed interface {
	Migrations() embed.FS
}
