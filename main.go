package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/getlantern/systray"
	"github.com/labstack/echo"
	"github.com/nightlyone/lockfile"
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
	lock, err := lockfile.New(filepath.Join(os.TempDir(), "go.pablo.lck"))

	if err != nil {
		msg := fmt.Sprintf("Cannot init lock. reason: %v", err)
		systray.ShowMessageBox("Lock error", msg)
		panic(err)
	}
	err = lock.TryLock()

	// Error handling is essential, as we only try to get the lock.
	if err != nil {
		systray.ShowMessageBox("Pablo", "Another process already running!")
		os.Exit(-1)
	}

	defer lock.Unlock()

	// Start server in another goroutine
	go startServer()

	// Should be called at the very beginning of main().
	systray.Run(onReady, func() {} /*onexit ignored*/)
}

func openURL() {
	url := fmt.Sprintf("http://localhost%s", serverAddr)
	open.Run(url)
}

func onReady() {

	systray.SetIcon(ICON_DATA)
	info := "pablo is running on port " + serverAddr
	systray.SetTooltip(info)

	menuOpen := systray.AddDefaultMenuItem("Open Pablo...", "Running on port "+serverAddr)
	systray.AddSeparator()
	menuUpdates := systray.AddMenuItem("Check for updates", "Check for updates")
	systray.AddSeparator()
	menuQuit := systray.AddMenuItem("Exit", "Quit the whole app")

	balloonClicked := systray.BalloonNotifyChan()

	systray.ShowBalloon("Pablo is running on port "+serverAddr, "Click to open your browser")

	go func() {
		for {
			select {
			case <-balloonClicked:
				openURL()
			case <-menuUpdates.ClickedCh:
				systray.ShowMessageBox("Pablo", "To be implemented")
			case <-menuOpen.ClickedCh:
				openURL()
			case <-menuQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}
