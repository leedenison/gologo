package timer

import "github.com/AllenDang/w32"
import "syscall"

var (
	moduser32 = syscall.NewLazyDLL("user32.dll")

    procSetTimer = moduser32.NewProc("SetTimer")
    procKillTimer = moduser32.NewProc("KillTimer")
)

func SetTimer(hwnd w32.HWND, nIDEvent uintptr, uElapse uint32, lpTimerProc uintptr) uintptr {
	ret, _, _ := procSetTimer.Call(
		uintptr(hwnd),
		nIDEvent,
		uintptr(uElapse),
		lpTimerProc)

	return ret
}

func KillTimer(hwnd w32.HWND, nIDEvent uintptr) bool {
	ret, _, _ := procKillTimer.Call(
		uintptr(hwnd),
		nIDEvent)

	return ret != 0
}