package assets

import "embed"

//go:embed "migrations/*.sql"
var MigrationsFiles embed.FS
