package main

import (
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	melody "gopkg.in/olahol/melody.v1"
)

func (s *Server) middlewares() {
	s.router.Pre(middleware.RemoveTrailingSlash())
	// s.router.Use(middleware.Logger())
	s.router.Use(middleware.Recover())
}

func (s *Server) routes() {
	assetHandler := http.FileServer(packr.NewBox("./static"))

	s.router.GET("/", echo.WrapHandler(assetHandler))
	s.router.GET("/___INTERNAL/static/*", echo.WrapHandler(http.StripPrefix("/___INTERNAL/static/", assetHandler)))
	s.router.GET("/___INTERNAL/ws", func(c echo.Context) error {
		s.ws.HandleRequest(c.Response(), c.Request())
		return nil
	})

	s.ws.HandleMessage(func(sess *melody.Session, msg []byte) {
		s.ws.Broadcast(msg)
	})

	s.router.POST("/*", s.handleUserPost)
}

func (s *Server) start(addr string, silent bool) {
	s.router.HideBanner = silent
	s.router.HidePort = silent
	s.router.Logger.Fatal(s.router.Start(addr))
}
