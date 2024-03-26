package describer

import (
	_ "embed"
)

//go:embed completions/describer.bash
var bashCompletion string

//go:embed completions/describer.zsh
var zshCompletion string

//go:embed completions/describer.ps1
var pwshCompletion string
