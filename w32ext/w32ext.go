package w32ext

import "github.com/AllenDang/w32"
import "syscall"
import "unsafe"

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

func KillTimer(wCtx *WindowContext, nIDEvent uintptr) bool {
	ret, _, _ := procKillTimer.Call(
		uintptr(wCtx.Window),
		nIDEvent)

	return ret != 0
}

func MakeIntResource(id uint16) *uint16 {
	return (*uint16)(unsafe.Pointer(uintptr(id)))
}

func RGB(r, g, b byte) uint32 {
	return (uint32(r) | uint32(g)<<8 | uint32(b)<<16)
}

type WindowContext struct {
	Window w32.HWND
	HDC w32.HDC
	previousPen w32.HGDIOBJ
	previousBrush w32.HGDIOBJ
	lBrushs map[Pen]*w32.LOGBRUSH
	hBrushs map[Pen]*w32.HBRUSH
	hPens map[Pen]*w32.HPEN
}

type Pen struct {
	Color uint32
}

func GetClientRect(wCtx *WindowContext) *w32.RECT {
	return w32.GetClientRect(wCtx.Window)
}

func SelectPen(wCtx *WindowContext, pen *Pen) {
	if wCtx.lBrushs == nil {
		wCtx.lBrushs = make(map[Pen]*w32.LOGBRUSH)
	}

	if wCtx.hBrushs == nil {
		wCtx.hBrushs = make(map[Pen]*w32.HBRUSH)
	}

	if wCtx.hPens == nil {
		wCtx.hPens = make(map[Pen]*w32.HPEN)
	}

	if wCtx.hPens[*pen] == nil {
	  // Create brush and pen
	  lBrush := w32.LOGBRUSH { LbStyle: w32.BS_SOLID, LbColor: w32.COLORREF(pen.Color) }
	  hBrush := w32.CreateBrushIndirect(&lBrush)
	  hPen := w32.ExtCreatePen(w32.PS_COSMETIC|w32.PS_SOLID, 1, &lBrush, 0, nil)
	  wCtx.lBrushs[*pen] = &lBrush
	  wCtx.hBrushs[*pen] = &hBrush
	  wCtx.hPens[*pen] = &hPen
	}

	// Select pen and store previous
	tempPen := w32.SelectObject(wCtx.HDC, w32.HGDIOBJ(*wCtx.hPens[*pen]))
	tempBrush := w32.SelectObject(wCtx.HDC, w32.HGDIOBJ(*wCtx.hBrushs[*pen]))

	if wCtx.previousPen == 0 {
		wCtx.previousPen = tempPen
	}

	if wCtx.previousBrush == 0 {
		wCtx.previousBrush = tempBrush
	}
}

func ReleaseContext(wCtx *WindowContext) {
	// Delete objects
	if wCtx.hBrushs != nil {	
		for pen, hBrush := range wCtx.hBrushs {
			delete(wCtx.hBrushs, pen)
			w32.DeleteObject(w32.HGDIOBJ(*hBrush))
		}
	}

	if wCtx.hPens != nil {	
		for pen, hPen := range wCtx.hPens {
			delete(wCtx.hPens, pen)
			w32.DeleteObject(w32.HGDIOBJ(*hPen))
		}
	}

	// Reselect previous pen
	if wCtx.previousPen == 0 {
		w32.SelectObject(wCtx.HDC, wCtx.previousPen)
	}

	if wCtx.previousBrush == 0 {
		w32.SelectObject(wCtx.HDC, wCtx.previousBrush)
	}
}

func DrawEllipse(wCtx *WindowContext, pen *Pen, l, t, r, b int32) {
	SelectPen(wCtx, pen)

	w32.Ellipse(wCtx.HDC, int(l), int(t), int(r), int(b))

	ReleaseContext(wCtx)
}

func DrawRectangle(wCtx *WindowContext, pen *Pen, l, t, r, b int32) {
	SelectPen(wCtx, pen)

	w32.Rectangle(wCtx.HDC, int(l), int(t), int(r), int(b))

	ReleaseContext(wCtx)
}

func ClearRect(wCtx *WindowContext, l, t, r, b int32) {
	rect := w32.RECT{Left: l, Top: t, Right: r, Bottom: b}
	w32.InvalidateRect(wCtx.Window, &rect, true)
}

func ReleaseDC(wCtx *WindowContext) {
	w32.ReleaseDC(wCtx.Window, wCtx.HDC)
}
