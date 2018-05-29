package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/kr/pretty"
	"github.com/labstack/echo"
)

func (s *Server) handleUserPost(c echo.Context) error {
	uri := c.Request().RequestURI
	// method := c.Request().Method
	// remote := c.Request().RemoteAddr

	// json_map := make(map[string]interface{})
	// json.NewDecoder(c.Request().Body).Decode(&json_map)
	body, _ := ioutil.ReadAll(c.Request().Body)
	var fields interface{}
	err := json.Unmarshal(body, &fields)
	if err != nil {
		panic(err)
	}

	var stream *Stream
	var ok bool

	if stream, ok = s.streams[uri]; ok {
		fmt.Printf("PATH KNOWN: %s\n", uri)
	} else {
		fmt.Printf("NEW PATH: %s\n", uri)
		s.streams[uri] = &Stream{
			varMap: make(map[string]string),
		}
		// Send to user?
		// O que acontece quando user abre a pagina novamente? Precisa buscar todos paths no open.
		// ou... Enviamos novamente... a cada post?...

	}

	msgBody := ""
	m := fields.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		// case string:
		// 	msgBody += fmt.Sprintf("\n%s is string %s", k, vv)
		// case bool:
		// 	msgBody += fmt.Sprintf("\n%s is boolean %v", k, vv)
		case float64:
			if _, ok := s.streams[uri].varMap[k]; !ok {
				fmt.Printf(" NEW VAR %s\n", k)
				if s.streams[uri].varMap == nil {
					s.streams[uri].varMap = make(map[string]string)
				}
				s.streams[uri].varMap[k] = "float64"
			}
			msgBody += fmt.Sprintf("\n%s is float64 %f", k, vv)
		default:
			msgBody += fmt.Sprintf("\n%s will be ignored", k)
		}
	}

	// msg := fmt.Sprintf("%s %s Body: %s", method, uri, msgBody)
	// fmt.Println(msg)

	// TODO:
	// Aqui entao pegar essa nova colecao, que pode ou nao ja estar salva na nossa memoria, e
	// adicionar as novas variaveis que achamos, ou modificar o tipo se for o caso.
	// O usuario por sua vez ja vai ver um popup de grafico aonde pode selecionar colecao e variaveis.
	// Com alguma opcao de converter o valor antes de mostrar (tipo unix nano seconds).
	//

	// TODO: send data to user for charting, if there is a chart already configured for this stream
	for _, chart := range s.charts {

		// If this POST is for a chart the user has configured.
		//if uri == chart.stream.path {
		if stream == chart.stream {
			fmt.Printf("Found configured chart %s\n", chart.stream.path)
		}

	}

	pretty.Print(s.streams)
	pretty.Print(s.charts)

	// s.ws.Broadcast([]byte(msg))
	return nil
}
