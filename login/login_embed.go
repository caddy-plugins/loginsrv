package login

import (
	_ "embed"
)

//go:embed assets/bootstrap.min.css
var bootstrapCSS string

//go:embed assets/bootstrap-social.min.css
var bsSocialCSS string

//go:embed assets/font-awesome.css
var fontAwesomeCSS string

//go:embed template/partial.html
var partials string

//go:embed template/layout.html
var layout string
