package focusedwindow

import (
	"syscall"
	"unicode/utf16"
)

// FocusedWindow stores relevant sys calls for determining a focused window
type FocusedWindow struct {
	user32 *syscall.LazyDLL
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getforegroundwindow
	procForeWin *syscall.LazyProc
	// https://learn.microsoft.com/en-us/windows/win32/api/winuser/nf-winuser-getwindowtextw
	procWinTxt *syscall.LazyProc
}

// NewFocusedWindow returns a focused window using user32
func NewFocusedWindow() FocusedWindow {
	user32 := syscall.NewLazyDLL("user32.dll")
	foreWin := user32.NewProc("GetForegroundWindow")
	winTxt := user32.NewProc("GetWindowTextW")

	return FocusedWindow{
		user32,
		foreWin,
		winTxt,
	}
}

// Title returns the title of a currently focused window or an error
func (s *FocusedWindow) Title() (string, error) {
	foreWinHandle, _, _ := s.procForeWin.Call()

	if foreWinHandle == 0 {
		return "", ErrNoFocusedWindow
	}

	winFocusedWindow := newText()
	FocusedWindowLen, _, _ := s.procWinTxt.Call(foreWinHandle, winFocusedWindow.pinter, winFocusedWindow.size)

	if FocusedWindowLen == 0 {
		return "", ErrNoFocusedWindow
	}

	FocusedWindow := string(utf16.Decode(winFocusedWindow.buffer[:FocusedWindowLen]))
	return FocusedWindow, nil
}
