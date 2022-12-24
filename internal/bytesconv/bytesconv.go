// Copyright 2020 Gin Core Team. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package bytesconv

import (
	"unsafe"
)

// StringToBytes 将字符串转换为字节切片片，无需内存分配
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// BytesToString 将字节片转换为字符串，无需内存分配
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
