package resources

import "embed"

//go:embed data/*.txt
var Dictionaries embed.FS
