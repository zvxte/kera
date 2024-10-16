package database

import "embed"

//go:embed migrations/*.sql
var assets embed.FS
