package login

import (
	"embed"
)

var bootstrapCSS string

var bsSocialCSS string

var fontAwesomeCSS string

//go:embed template/partial.html
var partials string

//go:embed template/layout.html
var layout string

//go:embed assets
var assetsFS embed.FS
