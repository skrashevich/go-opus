package libc

import "unsafe"

// VaList stores args at p, 8 bytes per argument, and returns p. It supports
// the argument types lib.go passes.
func VaList(p uintptr, args ...any) uintptr {
	r := p
	for _, v := range args {
		switch x := v.(type) {
		case int32:
			*(*int64)(unsafe.Pointer(p)) = int64(x)
		case int64:
			*(*int64)(unsafe.Pointer(p)) = x
		case uint32:
			*(*uint64)(unsafe.Pointer(p)) = uint64(x)
		case uint64:
			*(*uint64)(unsafe.Pointer(p)) = x
		case uintptr:
			*(*uintptr)(unsafe.Pointer(p)) = x
		default:
			panic("libc.VaList: unsupported argument type")
		}
		p += 8
	}
	return r
}

func VaInt32(app *uintptr) int32 {
	ap := *app
	if ap == 0 {
		return 0
	}
	*app = ap + 8
	return int32(*(*int64)(unsafe.Pointer(ap)))
}

func VaUintptr(app *uintptr) uintptr {
	ap := *app
	if ap == 0 {
		return 0
	}
	*app = ap + 8
	return *(*uintptr)(unsafe.Pointer(ap))
}
