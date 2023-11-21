package assets

import "embed"

//go:embed "migrations" "emails"
var EmbeddedFiles embed.FS
