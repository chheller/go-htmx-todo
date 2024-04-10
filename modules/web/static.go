package web

import (
	"embed"
	"io/fs"
)

var (
	//go:embed all:static
	StaticWebFiles embed.FS
)

func GetStaticWebAssets() (fs.FS, error) {
	return fs.Sub(StaticWebFiles, "static")
}
