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

import "net/http"

type CallOption interface {
	before(*callInfo) error
	after(*callInfo, *csAttempt)
}

type csAttempt struct {
	response *http.Response
}

// callInfo 调用信息
type callInfo struct {
	contentType  string
	operation    string
	pathTemplate string
}

// newDefaultCallInfo 默认调用信息
func newDefaultCallInfo(path string) callInfo {
	return callInfo{
		contentType:  "application/json",
		operation:    path,
		pathTemplate: path,
	}
}

type EmptyCallOption struct{}

func (EmptyCallOption) before(*callInfo) error      { return nil }
func (EmptyCallOption) after(*callInfo, *csAttempt) {}

// Operation is serviceMethod call option
func Operation(operation string) CallOption {
	return OperationCallOption{Operation: operation}
}

// OperationCallOption is set ServiceMethod for client call
type OperationCallOption struct {
	EmptyCallOption
	Operation string
}

func (o OperationCallOption) before(c *callInfo) error {
	c.operation = o.Operation
	return nil
}

// PathTemplate is http path template
func PathTemplate(pattern string) CallOption {
	return PathTemplateCallOption{Pattern: pattern}
}

// PathTemplateCallOption is set path template for client call
type PathTemplateCallOption struct {
	EmptyCallOption
	Pattern string
}

func (o PathTemplateCallOption) before(c *callInfo) error {
	c.pathTemplate = o.Pattern
	return nil
}
