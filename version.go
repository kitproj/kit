package main

import (
	_ "embed"
)

//go:generate sh -c "git describe --tags > tag"
//go:embed tag
var tag string
