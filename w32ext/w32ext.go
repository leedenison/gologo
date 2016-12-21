package w32ext

import "github.com/AllenDang/w32"
import "syscall"
import "unsafe"

var (
    moduser32 = syscall.NewLazyDLL("user32.dll")
    gdi32 = syscall.NewLazyDLL("gdi32.dll")

    procSetTimer = moduser32.NewProc("SetTimer")
    procKillTimer = moduser32.NewProc("KillTimer")
    procCreateCompatibleBitmap = gdi32.NewProc("CreateCompatibleBitmap")
    procGetCurrentObject = gdi32.NewProc("GetCurrentObject")
)

const (
    OBJ_BITMAP = 7
)

func SetTimer(
        hwnd w32.HWND,
        nIDEvent uint32,
        uElapse uint32,
        lpTimerProc uintptr) uintptr {
    ret, _, _ := procSetTimer.Call(
        uintptr(hwnd),
        uintptr(nIDEvent),
        uintptr(uElapse),
        lpTimerProc)

    return ret
}

func KillTimer(hwnd w32.HWND, nIDEvent uint32) bool {
    ret, _, _ := procKillTimer.Call(
        uintptr(hwnd),
        uintptr(nIDEvent))

    return ret != 0
}

func CreateCompatibleBitmap(hdc w32.HDC, width, height int32) w32.HBITMAP {
    ret, _, _ := procCreateCompatibleBitmap.Call(
        uintptr(hdc),
        uintptr(width),
        uintptr(height))

    return w32.HBITMAP(ret)
}

func GetCurrentObject(hdc w32.HDC, uObjectType uint32) w32.HGDIOBJ {
    ret, _, _ := procGetCurrentObject.Call(
        uintptr(hdc),
        uintptr(uObjectType))

    return w32.HGDIOBJ(ret)
}

func MakeIntResource(id uint16) *uint16 {
    return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func RGB(r, g, b byte) uint32 {
    return (uint32(r) | uint32(g)<<8 | uint32(b)<<16)
}

type Event struct {
    Id uintptr
    Time uint16
    WParam uintptr
    LParam uintptr
}

