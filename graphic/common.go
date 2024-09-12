package graphic

import (
	"unsafe"
)

// map a cordinate (x, y) in the infinite plane to a rectangle size rectW * rectH at (rectX, rectY)
// the coordinate is also mapped to the rectangle's native resolution of [0, rectNativeW) * [0, rectNativeH)
// this function return (-1, -1) if (x, y) is outside
func MapCoordinateIfInside(x, y, rectX, rectY, rectW, rectH, rectNativeW, rectNativeH int) (int, int) {
	if (x < rectX) || (y < rectY) || (x >= rectX+rectW) || (y >= rectY+rectH) {
		return -1, -1
	}
	return (x - rectX) * rectNativeW / rectW, (y - rectY) * rectNativeH / rectH
}

func StringFromUTF32(utf32 *uint) string {
	ptr := uintptr(unsafe.Pointer(utf32))
	s := ""
	for {
		r := *((*rune)(unsafe.Pointer(ptr)))
		ptr += 4
		if r == 0 {
			break
		}
		s += string(r)
	}
	return s
}
