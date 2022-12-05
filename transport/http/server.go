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

package http

import "github.com/go-ceres/ceres/server"

type HandleFunc func(Context) error

const SupportPackageIsVersion1 = true

type Server interface {
	server.Server
	Handler(method, path string, h HandleFunc, filters ...interface{})
	GET(path string, h HandleFunc, filters ...interface{})
	POST(path string, h HandleFunc, filters ...interface{})
	HEAD(path string, h HandleFunc, filters ...interface{})
	PUT(path string, h HandleFunc, filters ...interface{})
	PATCH(path string, h HandleFunc, filters ...interface{})
	DELETE(path string, h HandleFunc, filters ...interface{})
	CONNECT(path string, h HandleFunc, filters ...interface{})
	OPTIONS(path string, h HandleFunc, filters ...interface{})
	TRACE(path string, h HandleFunc, filters ...interface{})
}
