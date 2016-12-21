package gologo

import (
    "github.com/AllenDang/w32"
    "github.com/leedenison/gologo/w32ext"
    "syscall"
    "unsafe"
)

const SIXTY_HZ_IN_MILLIS = 16

const GOLOGO_MAIN_WIN = "GOLOGO_MAIN"

const WIN_SIZE_X = 1024
const WIN_SIZE_Y = 768

const TIMER_ID = iota

var EventHandlers = map[uint32]func(w32.HWND, *w32ext.Event) {}

func TimerProc(
        hwnd w32.HWND,
        uMsg uint32,
        idEvent uintptr,
        dwTime uint16) uintptr {
    ev := &w32ext.Event { Id: idEvent, Time: dwTime }

    EventHandlers[w32.WM_TIMER](hwnd, ev)

    return 0
}

func WndProc(hwnd w32.HWND, msg uint32, wParam, lParam uintptr) uintptr {
    ev := &w32ext.Event { WParam: wParam, LParam: lParam }

    switch msg {
    case w32.WM_DESTROY:
        // 0 = WM_QUIT
        w32.PostQuitMessage(0)
    case w32.WM_ERASEBKGND:
        return 1 // Custom handling
    case w32.WM_SIZE, w32.WM_PAINT:
        // On initial paint
        var ps w32.PAINTSTRUCT

        w32.BeginPaint(hwnd, &ps)
        EventHandlers[msg](hwnd, ev)
        w32.EndPaint(hwnd, &ps)
    default:
        return w32.DefWindowProc(hwnd, msg, wParam, lParam)
    }
    return 0
}

func CreateWindowClass(
        app w32.HINSTANCE,
        className string) w32.WNDCLASSEX {
    var wcex w32.WNDCLASSEX

    // Registered class name of the window.
    lpszClassName := syscall.StringToUTF16Ptr(className)

    // Size of the window object.
    wcex.Size = uint32(unsafe.Sizeof(wcex))

    wcex.Style = w32.CS_HREDRAW | w32.CS_VREDRAW
    // Application loop handler procedure.
    wcex.WndProc = syscall.NewCallback(WndProc)
    // Additional bytes to allocate for the window class struct.
    wcex.ClsExtra = 0
    // Additional bytes to allocation for the window instance struct.
    // If an application uses WNDCLASS to register a dialog box created
    // by using the CLASS directive in the resource file, it must set this
    // member to DLGWINDOWEXTRA.
    wcex.WndExtra = 0
    // A handle to the instance that contains the window procedure for the
    // class.
    wcex.Instance = app

    // Use default IDI_APPLICATION icon.
    wcex.Icon = w32.LoadIcon(
        app,
        w32ext.MakeIntResource(w32.IDI_APPLICATION))

    // Use default IDC_ARROW mouse cursor.
    wcex.Cursor = w32.LoadCursor(0, w32ext.MakeIntResource(w32.IDC_ARROW))

    // Assign HBRUSH to background using the standard window color
    wcex.Background = w32.COLOR_WINDOW + 11

    // Assign name of menu resource
    wcex.MenuName = nil

    wcex.ClassName = lpszClassName
    wcex.IconSm = w32.LoadIcon(
        app,
        w32ext.MakeIntResource(w32.IDI_APPLICATION))

    // Make this window available to other controls
    w32.RegisterClassEx(&wcex)

    return wcex
}

func CreateWindowInstance(
        app w32.HINSTANCE,
        className string,
        title string) w32.HWND {
        hwnd := w32.CreateWindowEx(0, syscall.StringToUTF16Ptr(className),
        syscall.StringToUTF16Ptr(title),
        w32.WS_OVERLAPPEDWINDOW|w32.WS_VISIBLE,
        w32.CW_USEDEFAULT, w32.CW_USEDEFAULT, WIN_SIZE_X, WIN_SIZE_Y, 0, 0,
        app, nil)

    return hwnd
}

func SetTimer(
        hwnd w32.HWND,
        nIDEvent uint32,
        uElapse uint32,
        timerCallback func(w32.HWND, *w32ext.Event)) {
    EventHandlers[w32.WM_TIMER] = timerCallback
    w32ext.SetTimer(hwnd, nIDEvent, uElapse, syscall.NewCallback(TimerProc))
}

func WindowsTick(hwnd w32.HWND, ev *w32ext.Event) {
    // TODO: Need to mutex this so we don't enter twice
    PhysicsTick(hwnd)
    PaintTick(hwnd, ev)
}

func Run(title string) {
    app := w32.GetModuleHandle("")

    EventHandlers[w32.WM_SIZE] = OnSize
    EventHandlers[w32.WM_PAINT] = OnPaint

    CreateWindowClass(app, GOLOGO_MAIN_WIN)
    hwnd := CreateWindowInstance(app, GOLOGO_MAIN_WIN, title)
    CreateBuffer(hwnd)
    CreateRenderers(hwnd)
    UpdateWindowEdge(hwnd)
    SetTimer(hwnd, TIMER_ID, SIXTY_HZ_IN_MILLIS, WindowsTick)

    var msg w32.MSG
    for {
        // 0, 0, 0 = retrive all messages from all sources
        if w32.GetMessage(&msg, 0, 0, 0) == 0 {
            break
        }
        w32.TranslateMessage(&msg)
        w32.DispatchMessage(&msg)
    }

    ReleaseRenderers(hwnd)
    ReleaseBuffer()

    return
}
