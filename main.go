package main

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/pkg/browser"
	"gopkg.in/olahol/melody.v1"
)

const serverAddr = ":5000"

// Server struct to share required data for HTTP handlers
type Server struct {
	ws     *melody.Melody
	router *echo.Echo
}

func main() {
	s := Server{
		router: echo.New(),
		ws:     melody.New(),
	}

	s.setup()

	url := fmt.Sprintf("http://localhost%s", serverAddr)
	browser.OpenURL(url)

	s.start(serverAddr)
}
