//    Copyright 2022. ceres
//    Author https://github.com/go-ceres/ceres
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package config

type Option func(*Options)

// Options 参数信息
type Options struct {
	sources  []Source
	decoder  Decoder
	resolver resolver
}

// WithDecoder 设置数据解码器
func WithDecoder(decoder Decoder) Option {
	return func(options *Options) {
		options.decoder = decoder
	}
}

// WithSource 设置数据解析器
func WithSource(sources ...Source) Option {
	return func(options *Options) {
		options.sources = append(options.sources, sources...)
	}
}

// WithResolver 设置数据解析器
func WithResolver(resolver resolver) Option {
	return func(options *Options) {
		options.resolver = resolver
	}
}
