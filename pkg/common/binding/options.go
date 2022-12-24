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

package binding

import (
	"github.com/andeya/goutil"
	"strings"
)

var (
	defaultTagPath   = "path"
	defaultTagQuery  = "query"
	defaultTagHeader = "header"
	defaultTagCookie = "cookie"
	defaultTagForm   = "form"
	tagProtobuf      = "protobuf"
	tagJSON          = "json"
)

type Option func(o *Options)

type Options struct {
	// LooseZeroMode if set to true,
	// the empty string request parameter is bound to the zero value of parameter.
	// NOTE: Suitable for these parameter types: query/header/cookie/form .
	LooseZeroMode bool
	// PathParam use 'path' by default when empty
	PathParam string
	// Query use 'query' by default when empty
	Query string
	// Header use 'header' by default when empty
	Header string
	// Cookie use 'cookie' by default when empty
	Cookie string
	// RawBody use 'raw' by default when empty
	RawBody string
	// FormBody use 'form' by default when empty
	FormBody string
	// Validator use 'vd' by default when empty
	Validator string
	// protobufBody use 'protobuf' by default when empty
	protobufBody string
	// jsonBody use 'json' by default when empty
	jsonBody string
	// defaultVal use 'default' by default when empty
	defaultVal  string
	pathReplace PathReplace
	list        []string
}

func DefaultOptions() *Options {
	options := new(Options)
	options.list = []string{
		goutil.InitAndGetString(&options.PathParam, defaultTagPath),
		goutil.InitAndGetString(&options.Query, defaultTagQuery),
		goutil.InitAndGetString(&options.Header, defaultTagHeader),
		goutil.InitAndGetString(&options.Cookie, defaultTagCookie),
		goutil.InitAndGetString(&options.FormBody, defaultTagForm),
		goutil.InitAndGetString(&options.protobufBody, tagProtobuf),
		goutil.InitAndGetString(&options.jsonBody, tagJSON),
	}
	options.pathReplace = func(name, value string, path string) string {
		return strings.Replace(path, ":"+name, value, -1)
	}
	return options
}
