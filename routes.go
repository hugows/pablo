package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func (s *Server) setup() {
	// Middlewares
	s.router.Pre(middleware.RemoveTrailingSlash())
	s.router.Use(middleware.Recover())

	// Packr is used so we can deliver a single binary to the end user
	assetHandler := http.FileServer(packr.NewBox("./static"))

	// Serve index.html
	s.router.GET("/", echo.WrapHandler(assetHandler))

	// Serve js/css/...
	s.router.GET("/___INTERNAL/static/*", echo.WrapHandler(http.StripPrefix("/___INTERNAL/static/", assetHandler)))

	// Start a websocket connection
	s.router.GET("/___INTERNAL/ws", func(c echo.Context) error {
		s.ws.HandleRequest(c.Response(), c.Request())
		return nil
	})

	// Receives POST data, filter numeric fields, adds special fields
	// (server timestamp, original URI) and broadcasts to the connected
	// client on the browser.
	s.router.POST("/*", func(c echo.Context) error {
		uri := c.Request().RequestURI

		body, _ := ioutil.ReadAll(c.Request().Body)
		var input interface{}

		err := json.Unmarshal(body, &input)
		if err != nil {
			return err
		}

		minput := input.(map[string]interface{})
		output := make(map[string]interface{})

		for k, v := range minput {
			switch vv := v.(type) {
			case float64:
				output[k] = vv
			}
		}

		output["_timestamp"] = time.Now()
		output["_uri"] = uri

		msg, err := json.Marshal(output)

		if err != nil {
			return err
		}

		s.ws.Broadcast(msg)
		return nil
	})
}

func (s *Server) start(addr string) {
	s.router.HideBanner = true
	s.router.Logger.Fatal(s.router.Start(addr))
}
