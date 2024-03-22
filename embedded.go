package tbdex

import "embed"

// EmbeddedFiles is used to access files
//
//go:embed  tbdex/hosted/*
var EmbeddedFiles embed.FS
