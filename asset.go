// +build !packr

package main

import (
	"net/http"

	"github.com/gobuffalo/packr"
)

func getAssetHandler() http.Handler {
	fs := packr.NewBox("./static")
	return http.FileServer(fs)
}
