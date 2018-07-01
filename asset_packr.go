// +build packr

package main

import "net/http"

func getAssetHandler() http.Handler {
	fs := http.Dir("./static")
	return http.FileServer(fs)
}
