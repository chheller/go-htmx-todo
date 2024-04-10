package web

import (
	"embed"
	"io/fs"
)

var (
	//go:embed static/**/*
	StaticWebFiles embed.FS
)

func GetStaticWebAssets() (fs.FS, error) {
	return fs.Sub(StaticWebFiles, "static")
}
