package main

import (
	"key-changer/src"
)

func main() {
	src.InitConfig()
	go src.InitTray()

	src.InitKeyChanger()
}
