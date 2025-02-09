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
	systray.SetTooltip("Making one key press act as another.")

	open := systray.AddMenuItem("Config", "Edit the configuration file.")
	reload := systray.AddMenuItem("Reload", "Reload the configuration.")
	systray.AddSeparator()
	keyPressItem = systray.AddMenuItemCheckbox("View key code", "Click to press a key and view its code.", false)
	systray.AddSeparator()
	exit := systray.AddMenuItem("Close", "Close Key-Changer.")

	go func() {
		for {
			select {
			case <-open.ClickedCh:
				cmd := exec.Command(KeyChangerConfig.Program, configFile)
				cmd.Stdout = nil
				cmd.Stderr = nil
				err := cmd.Start()
				if err != nil {
					log.Println("Error opening the configuration file with \""+KeyChangerConfig.Program+"\":", err)
					return
				}
			case <-keyPressItem.ClickedCh:
				if keyPressItem.Checked() {
					keyPressItem.Uncheck()
					keyPressItem.SetTitle("View key code")
				} else {
					keyPressItem.SetTitle("Press-Key")
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
