package main

import (
	"github.com/labstack/echo"
	"gopkg.in/olahol/melody.v1"
)

// Stream represents a data source - each unique URI posted by the user becomes a stream.
// The user may choose to hide any streams (accidental POST).
type Stream struct {
	path   string
	varMap map[string]string
	hidden bool
}

// Chart are the charts being displayed on the current page.
// The user creates those charts after seeing the streams pop-up on the page.
type Chart struct {
	// Stream (currently the "unique URL") to use as a data source for this chart
	stream *Stream

	// Extract from JSON to plot on the Y axis
	keyY string

	// Extract from JSON to plot on the X axis
	keyX string

	// Instead of getting an X value from the JSON struct, the user can also use the time
	// when the POST was received
	UseServerTimeForX bool
}

// Server struct to share required data for HTTP handlers
type Server struct {
	ws      *melody.Melody
	router  *echo.Echo
	streams map[string]*Stream
	charts  []*Chart
}

func main() {
	s := Server{
		router:  echo.New(),
		ws:      melody.New(),
		streams: make(map[string]*Stream),
		charts:  make([]*Chart, 0),
	}

	// DEBUG

	s.streams["/api/update"] = &Stream{
		path: "/api/update",
	}

	s.charts = append(s.charts, &Chart{
		stream: s.streams["/api/update"],
		keyX:   "temp",
		keyY:   "press",
	})

	s.middlewares()
	s.routes()

	// url := fmt.Sprintf("http://localhost%s", addr)
	// browser.OpenURL(url)

	s.start(":5000", false)
}
