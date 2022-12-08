// Copyright 2022. ceres
// Author https://github.com/go-ceres/ceres
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package token

import "unsafe"

// contains 判断数组当中是否包含
func contains[T comparable](slice []T, obj T) bool {
	if len(slice) == 0 {
		return false
	}
	for _, t := range slice {
		if t == obj {
			return true
		}
	}
	return false
}

// stringToBytes converts string to byte slice without a memory allocation.
func stringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// bytesToString converts byte slice to string without a memory allocation.
func bytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
