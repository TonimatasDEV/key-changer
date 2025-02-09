package src

import (
	"golang.org/x/sys/windows"
	"log"
	"strconv"
	"syscall"
	"unsafe"
)

var (
	user32     = windows.NewLazySystemDLL("user32.dll")
	setHook    = user32.NewProc("SetWindowsHookExW")
	callNext   = user32.NewProc("CallNextHookEx")
	getMsg     = user32.NewProc("GetMessageW")
	keybdEvent = user32.NewProc("keybd_event")
	hook       windows.Handle
)

const (
	WhKeyboardLl = 13
	WmKeydown    = 0x0100
)

type KBDLLHOOKSTRUCT struct {
	VkCode, ScanCode, Flags, Time uint32
	DwExtraInfo                   uintptr
}

func keyboardHook(nCode int, wParam, lParam uintptr) uintptr {
	if nCode >= 0 && wParam == WmKeydown {
		kb := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))

		if keyPressItem.Checked() {
			keyPressItem.SetTitle("Key: " + strconv.Itoa(int(kb.VkCode)))
			return 1
		}

		for _, key := range KeyChangerConfig.Keys {
			if kb.VkCode == key.KeyFrom {
				keybdEvent.Call(key.KeyTo, 0, 0, 0) // KeyDown
				keybdEvent.Call(key.KeyTo, 0, 2, 0) // KeyUp
				return 1
			}
		}
	}
	ret, _, _ := callNext.Call(0, uintptr(nCode), wParam, lParam)
	return ret
}

func InitKeyChanger() {
	ret, _, _ := setHook.Call(WhKeyboardLl, syscall.NewCallback(keyboardHook), 0, 0)
	hook = windows.Handle(ret)

	if hook == 0 {
		log.Fatalln("Error initializing keyboard hook.")
	}

	log.Println("Keyboard hook enabled...")
	var msg struct{}
	for {
		getMsg.Call(uintptr(unsafe.Pointer(&msg)), 0, 0, 0)
	}
}
