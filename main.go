package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gobuffalo/packr"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gopkg.in/olahol/melody.v1"
)

func main() {
	e := echo.New()
	m := melody.New()

	e.Pre(middleware.RemoveTrailingSlash())
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// to pack our assets (index.html, favicon..)
	assetHandler := http.FileServer(packr.NewBox("./static"))

	// serve index.html
	e.GET("/", echo.WrapHandler(assetHandler))

	// serve other files
	e.GET("/___INTERNAL/static/*", echo.WrapHandler(http.StripPrefix("/___INTERNAL/static/", assetHandler)))

	// melody starts connection using
	e.GET("/___INTERNAL/ws", func(c echo.Context) error {
		m.HandleRequest(c.Response(), c.Request())
		return nil
	})

	e.POST("/___INTERNAL/*", func(c echo.Context) error {
		// nop
		return nil
	})

	// send message to connected clients
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		m.Broadcast(msg)
	})

	e.POST("/*", func(c echo.Context) error {
		uri := c.Request().RequestURI
		method := c.Request().Method
		// remote := c.Request().RemoteAddr

		// json_map := make(map[string]interface{})
		// json.NewDecoder(c.Request().Body).Decode(&json_map)
		body, _ := ioutil.ReadAll(c.Request().Body)
		var fields interface{}
		err := json.Unmarshal(body, &fields)
		if err != nil {
			panic(err)
		}

		msgBody := ""
		mfields := fields.(map[string]interface{})

		for k, v := range mfields {
			switch vv := v.(type) {
			case string:
				msgBody += fmt.Sprintf("\n%s is string %s", k, vv)
			case bool:
				msgBody += fmt.Sprintf("\n%s is boolean %v", k, vv)
			case float64:
				msgBody += fmt.Sprintf("\n%s is float64 %f", k, vv)
			default:
				msgBody += fmt.Sprintf("\n%s is of a type I don't know how to handle", k)
			}
		}

		msg := fmt.Sprintf("%s %s Body: %s", method, uri, msgBody)
		fmt.Println(msg)

		// TODO:
		// Aqui entao pegar essa nova colecao, que pode ou nao ja estar salva na nossa memoria, e
		// adicionar as novas variaveis que achamos, ou modificar o tipo se for o caso.
		// O usuario por sua vez ja vai ver um popup de grafico aonde pode selecionar colecao e variaveis.
		// Com alguma opcao de converter o valor antes de mostrar (tipo unix nano seconds).
		//

		// m.Broadcast([]byte(msg))
		return nil
	})

	// start
	e.HideBanner = true
	e.HidePort = true

	addr := ":5000"
	// url := fmt.Sprintf("http://localhost%s", addr)
	// browser.OpenURL(url)
	//e.Logger.Fatal(e.Start(addr))
	e.Start(addr)
}
