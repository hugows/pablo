package main

import (
	"fmt"

	"github.com/getlantern/systray"
	"github.com/labstack/echo"
	"github.com/skratchdot/open-golang/open"
	melody "gopkg.in/olahol/melody.v1"
)

const serverAddr = ":5000"

// Server struct to share required data for HTTP handlers
type Server struct {
	ws     *melody.Melody
	router *echo.Echo
}

func startServer() {
	e := echo.New()
	e.Debug = true

	s := Server{
		router: e,
		ws:     melody.New(),
	}

	s.setup()
	s.start(serverAddr)
}

func main() {
	onExit := func() {
		fmt.Println("Starting onExit")
		// now := time.Now()
		// ioutil.WriteFile(fmt.Sprintf(`on_exit_%d.txt`, now.UnixNano()), []byte(now.String()), 0644)
		fmt.Println("Finished onExit")
	}

	// Start server in another goroutine
	go startServer()

	// Should be called at the very beginning of main().
	systray.Run(onReady, onExit)
}

func openURL() {
	url := fmt.Sprintf("http://localhost%s", serverAddr)
	open.Run(url)
}

func onReady() {

	systray.SetIcon(ICON_DATA)
	info := "pablo is running on port " + serverAddr
	systray.SetTitle(info)
	systray.SetTooltip(info)

	menuOpen := systray.AddMenuItem("Open Pablo...", "Running on port "+serverAddr)
	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Exit", "Quit the whole app")

	go func() {
		for {
			select {
			case <-menuOpen.ClickedCh:
				openURL()
			case <-menuQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}
