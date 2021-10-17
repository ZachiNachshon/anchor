package resources

import "embed"

//go:embed version.txt
var VersionFile embed.FS

//go:embed manifest.json
var ManifestFile embed.FS
