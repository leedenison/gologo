package gologo

import (
    "github.com/AllenDang/w32"
    "github.com/leedenison/gologo/w32ext"
)

/*
struct BITMAP {
  long           bmType;
  long           bmWidth;
  long           bmHeight;
  long           bmWidthBytes;
  unsigned short bmPlanes;
  unsigned short bmBitsPixel;
  void*          bmBits;
};
*/
import "C"

const PAINTER_BG = 0
const PAINTER_OBJ = 1

var windowHeight = int32(0)

var buffer = Buffer{}

var painters = map[int]*Painter {
    PAINTER_BG: &Painter{
        Type: OBJECT_SOLID,
        FGColor: w32ext.RGB(0, 0, 0),
        BGColor: w32ext.RGB(0, 0, 0),
    },
    PAINTER_OBJ: &Painter{
        Type: OBJECT_SOLID,
        FGColor: w32ext.RGB(0, 255, 0),
        BGColor: w32ext.RGB(0, 255, 0),
    },
}

type Buffer struct {
    HDC w32.HDC
    Bitmap w32.HBITMAP
    PreviousBitmap w32.HGDIOBJ
    Width int32
    Height int32
}

type Painter struct {
    Type int
    FGColor uint32
    BGColor uint32
    HPen w32.HPEN
    LBrush w32.LOGBRUSH
    HBrush w32.HBRUSH
}

type RenderState struct {
    Painter *Painter
    Pen w32.HPEN
    Brush w32.HBRUSH
}

func OnSize(hwnd w32.HWND, ev *w32ext.Event) {
    clientRect := w32.GetClientRect(hwnd)
    windowHeight = clientRect.Bottom
    OnPaint(hwnd, ev)
}

func OnPaint(hwnd w32.HWND, ev *w32ext.Event) {
    if buffer.HDC != 0 {
        ClearBuffer(hwnd, buffer.HDC, painters[PAINTER_BG])
        PaintObjects(hwnd, buffer.HDC)
        SwapBuffer(hwnd, buffer.HDC)
    }
}

func PaintTick(hwnd w32.HWND, ev *w32ext.Event) {
    OnPaint(hwnd, ev)
}

func PaintObjects(hwnd w32.HWND, hdc w32.HDC) {
    var previous RenderState

    for _, obj := range objects {
        obj.GetRenderer().Render(obj, hdc, &previous)
    }

    if previous.Pen != 0 {
        w32.SelectObject(hdc, w32.HGDIOBJ(previous.Pen))
    }

    if previous.Brush != 0 {
        w32.SelectObject(hdc, w32.HGDIOBJ(previous.Brush))
    }
}

func SelectPainterIfNeeded(
        hdc w32.HDC,
        Painter *Painter,
        previous *RenderState) {
    if Painter != previous.Painter {
        previous.Painter = Painter
        tempPen, tempBrush := SelectPainter(hdc, Painter)

        if previous.Pen == 0 {
            previous.Pen = tempPen
        }

        if previous.Brush == 0 {
            previous.Brush = tempBrush
        }
    }
}

func SelectPainter(hdc w32.HDC, Painter *Painter) (w32.HPEN, w32.HBRUSH) {
    // Select pen and store previous
    prevPen := w32.SelectObject(hdc, w32.HGDIOBJ(Painter.HPen))
    prevBrush := w32.SelectObject(hdc, w32.HGDIOBJ(Painter.HBrush))

    return w32.HPEN(prevPen), w32.HBRUSH(prevBrush)
}

func CreateBuffer(hwnd w32.HWND) {
    hdc := w32.GetDC(hwnd)
    clientRect := GetScreenRect(hwnd)

    buffer.HDC = w32.CreateCompatibleDC(hdc)
    buffer.Bitmap = w32ext.CreateCompatibleBitmap(hdc, clientRect.Right,
            clientRect.Bottom)

    buffer.PreviousBitmap = w32.SelectObject(buffer.HDC,
            w32.HGDIOBJ(buffer.Bitmap))
}

func ReleaseBuffer() {
    if buffer.HDC != 0 {
      w32.SelectObject(buffer.HDC, buffer.PreviousBitmap)
      w32.DeleteObject(w32.HGDIOBJ(buffer.Bitmap))
      w32.DeleteDC(buffer.HDC)
      buffer.Bitmap = 0
      buffer.HDC = 0
    }
}

func SwapBuffer(hwnd w32.HWND, hdc w32.HDC) {
    wHdc := w32.GetDC(hwnd)
    clientRect := w32.GetClientRect(hwnd)

    w32.BitBlt(
        wHdc,
        int(clientRect.Left),
        int(clientRect.Top),
        int(clientRect.Right),
        int(clientRect.Bottom),
        hdc,
        0,
        0,
        w32.SRCCOPY)

    w32.ReleaseDC(hwnd, wHdc)
}

func ClearBuffer(hwnd w32.HWND, hdc w32.HDC, Painter *Painter) {
    clientRect := w32.GetClientRect(hwnd)

    prevPen, prevBrush := SelectPainter(hdc, Painter)

    w32.Rectangle(
        hdc,
        int(clientRect.Left),
        int(clientRect.Top),
        int(clientRect.Right),
        int(clientRect.Bottom))

    if prevPen != 0 {
        w32.SelectObject(hdc, w32.HGDIOBJ(prevPen))
    }

    if prevBrush != 0 {
        w32.SelectObject(hdc, w32.HGDIOBJ(prevBrush))
    }
}

func GetScreenRect(hwnd w32.HWND) w32.RECT {
    hdc := w32.GetDC(hwnd)
    width := w32.GetDeviceCaps(hdc, w32.HORZRES)
    height := w32.GetDeviceCaps(hdc, w32.VERTRES)

    return w32.RECT {
        Left: 0,
        Right: int32(width),
        Top: 0,
        Bottom: int32(height),
    }
}

func CreatePainters(hwnd w32.HWND) {
    for _, Painter := range painters {
        // Create brush and pen
        Painter.LBrush = w32.LOGBRUSH {
            LbStyle: w32.BS_SOLID,
            LbColor: w32.COLORREF(Painter.BGColor),
        }
        Painter.HBrush = w32.CreateBrushIndirect(&Painter.LBrush)
        Painter.HPen = w32.ExtCreatePen(
            w32.PS_COSMETIC | w32.PS_SOLID,
            1,
            &Painter.LBrush,
            0,
            nil)
    }
}

func ReleasePainters(hwnd w32.HWND) {
    for _, Painter := range painters {
        if Painter.HBrush != 0 {
            w32.DeleteObject(w32.HGDIOBJ(Painter.HBrush))
            Painter.HBrush = 0
        }

        if Painter.HPen != 0 {
            w32.DeleteObject(w32.HGDIOBJ(Painter.HPen))
            Painter.HPen = 0
        }
    }
}