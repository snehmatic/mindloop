// Package web provides embedded static assets and templates
package web

import (
	"embed"
)

// WebFS contains the entire web directory embedded into the binary
//
//go:embed templates/* static/*
var WebFS embed.FS
