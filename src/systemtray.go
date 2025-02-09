package src

import (
	_ "embed"
	"github.com/getlantern/systray"
	"log"
	"os"
	"os/exec"
)

var (
	//go:embed icon.ico
	iconData     []byte
	keyPressItem *systray.MenuItem
)

func InitTray() {
	systray.Run(onReady, onExit)
}

func onReady() {
	systray.SetIcon(iconData)
	systray.SetTitle("Key-Changer")
	systray.SetTooltip("Change keyboard keys to other.")

	open := systray.AddMenuItem("Config", "Open config file.")
	reload := systray.AddMenuItem("Reload", "Reload config.")
	systray.AddSeparator()
	keyPressItem = systray.AddMenuItemCheckbox("Press-Key", "Click this and press a key to know code.", false)
	systray.AddSeparator()
	exit := systray.AddMenuItem("Close", "Close Key-Changer.")

	go func() {
		for {
			select {
			case <-open.ClickedCh:

				cmd := exec.Command("code", configFile)
				cmd.Stdout = nil
				cmd.Stderr = nil
				cmd.Start()
			case <-keyPressItem.ClickedCh:
				if keyPressItem.Checked() {
					keyPressItem.Uncheck()
				} else {
					keyPressItem.Check()
				}
			case <-reload.ClickedCh:
				reloadConfig()
				log.Println("Config reloaded successfully.")
			case <-exit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	log.Println("Stopping Key-Changer...")
	os.Exit(0)
}
