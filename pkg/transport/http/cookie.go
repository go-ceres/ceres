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

import (
	"github.com/valyala/fasthttp"
)

type CookieOption func(cookie *fasthttp.Cookie)

func CookieWithMaxAge(maxAge int) CookieOption {
	return func(cookie *fasthttp.Cookie) {
		cookie.SetMaxAge(maxAge)
	}
}

func CookieWithPath(path string) CookieOption {
	return func(cookie *fasthttp.Cookie) {
		cookie.SetPath(path)
	}
}

func CookieWithDomain(domain string) CookieOption {
	return func(cookie *fasthttp.Cookie) {
		cookie.SetDomain(domain)
	}
}

func CookieWithSecure(secure bool) CookieOption {
	return func(cookie *fasthttp.Cookie) {
		cookie.SetSecure(secure)
	}
}

func CookieWithHTTPOnly(httpOnly bool) CookieOption {
	return func(cookie *fasthttp.Cookie) {
		cookie.SetHTTPOnly(httpOnly)
	}
}

func SetSameSite(sameSite fasthttp.CookieSameSite) CookieOption {
	return func(cookie *fasthttp.Cookie) {
		cookie.SetSameSite(sameSite)
	}
}
