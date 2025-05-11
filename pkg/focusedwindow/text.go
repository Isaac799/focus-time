package focusedwindow

import "unsafe"

type text struct {
	buffer []uint16
	size   uintptr
	pinter uintptr
}

func newText() text {
	size := uintptr(256)
	buff := make([]uint16, size)
	pinter := uintptr(unsafe.Pointer(&buff[0]))

	return text{
		buff,
		size,
		pinter,
	}
}
