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

const RENDER_BG = 0
const RENDER_OBJ = 1

var buffer = Buffer{}

var renderers = map[int]*Renderer {
    RENDER_BG: &Renderer{
        Type: OBJECT_SOLID,
        FGColor: w32ext.RGB(0, 0, 0),
        BGColor: w32ext.RGB(0, 0, 0),
    },
    RENDER_OBJ: &Renderer{
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

type Renderer struct {
    Type int
    FGColor uint32
    BGColor uint32
    HPen w32.HPEN
    LBrush w32.LOGBRUSH
    HBrush w32.HBRUSH
}

func OnSize(hwnd w32.HWND, ev *w32ext.Event) {
    UpdateWindowEdge(hwnd)
    OnPaint(hwnd, ev)
}

func OnPaint(hwnd w32.HWND, ev *w32ext.Event) {
    if buffer.HDC != 0 {
        ClearBuffer(hwnd, buffer.HDC, renderers[RENDER_BG])
        PaintMovables(hwnd, buffer.HDC)
        SwapBuffer(hwnd, buffer.HDC)
    }
}

func PaintTick(hwnd w32.HWND, ev *w32ext.Event) {
    OnPaint(hwnd, ev)
}

func PaintMovables(hwnd w32.HWND, hdc w32.HDC) {
    var prevRenderer *Renderer
    var prevPen w32.HPEN
    var prevBrush w32.HBRUSH

    for _, obj := range objects {
        switch t := obj.(type) {
        case *Circle:
            if t.Renderer.Type != OBJECT_EMPTY {
                prevRenderer, prevPen, prevBrush =
                    SelectRendererIfNeeded(
                        hdc, 
                        t.Renderer,
                        prevRenderer,
                        prevPen,
                        prevBrush)

                w32.Ellipse(
                    hdc,
                    int(t.Center.x - t.Radius),
                    int(t.Center.y - t.Radius),
                    int(t.Center.x + t.Radius),
                    int(t.Center.y + t.Radius))
            }
        default:
            _ = t
        }
    }

    if prevPen != 0 {
        w32.SelectObject(hdc, w32.HGDIOBJ(prevPen))
    }

    if prevBrush != 0 {
        w32.SelectObject(hdc, w32.HGDIOBJ(prevBrush))
    }
}

func SelectRendererIfNeeded(
        hdc w32.HDC,
        renderer *Renderer,
        prevRenderer *Renderer,
        prevPen w32.HPEN,
        prevBrush w32.HBRUSH) (*Renderer, w32.HPEN, w32.HBRUSH) {
    if renderer != prevRenderer {
        prevRenderer = renderer
        tempPen, tempBrush := SelectRenderer(hdc, renderer)

        if prevPen == 0 {
            prevPen = tempPen
        }

        if prevBrush == 0 {
            prevBrush = tempBrush
        }
    }

    return prevRenderer, prevPen, prevBrush
}

func SelectRenderer(hdc w32.HDC, renderer *Renderer) (w32.HPEN, w32.HBRUSH) {
    // Select pen and store previous
    prevPen := w32.SelectObject(hdc, w32.HGDIOBJ(renderer.HPen))
    prevBrush := w32.SelectObject(hdc, w32.HGDIOBJ(renderer.HBrush))

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

func ClearBuffer(hwnd w32.HWND, hdc w32.HDC, renderer *Renderer) {
    clientRect := w32.GetClientRect(hwnd)

    prevPen, prevBrush := SelectRenderer(hdc, renderer)

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

func CreateRenderers(hwnd w32.HWND) {
    for _, renderer := range renderers {
        // Create brush and pen
        renderer.LBrush = w32.LOGBRUSH {
            LbStyle: w32.BS_SOLID,
            LbColor: w32.COLORREF(renderer.BGColor),
        }
        renderer.HBrush = w32.CreateBrushIndirect(&renderer.LBrush)
        renderer.HPen = w32.ExtCreatePen(
            w32.PS_COSMETIC | w32.PS_SOLID,
            1,
            &renderer.LBrush,
            0,
            nil)
    }
}

func ReleaseRenderers(hwnd w32.HWND) {
    for _, renderer := range renderers {
        if renderer.HBrush != 0 {
            w32.DeleteObject(w32.HGDIOBJ(renderer.HBrush))
            renderer.HBrush = 0
        }

        if renderer.HPen != 0 {
            w32.DeleteObject(w32.HGDIOBJ(renderer.HPen))
            renderer.HPen = 0
        }
    }
}