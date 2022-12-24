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

package i18n

import (
	"context"
	"github.com/go-ceres/ceres/pkg/common/codec"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"io/ioutil"
	"path/filepath"
)

type Option func(opts *options)

type options struct {
	languagePath    string
	format          string
	getLangHandler  func(ctx context.Context, def string) string
	defaultLanguage language.Tag
	acceptLanguage  []language.Tag
}

// buildBundle 构建bundle
func (o *options) buildBundle() *i18n.Bundle {
	// 创建bundle
	bundle := i18n.NewBundle(o.defaultLanguage)
	encoding := codec.LoadCodec(o.format)
	if encoding == nil {
		panic("get codec error：not set")
	}
	// 注册消息解码器
	bundle.RegisterUnmarshalFunc(o.format, encoding.Unmarshal)
	// 加载文件
	o.loadFile(bundle)
	// 返回bundle
	return bundle
}

func (o *options) buildLocalizers(bundle *i18n.Bundle) map[string]*i18n.Localizer {
	res := make(map[string]*i18n.Localizer, 0)
	for _, tag := range o.acceptLanguage {
		langStr := tag.String()
		res[langStr] = o.newLocalizer(bundle, langStr)
	}
	defaultLng := o.defaultLanguage.String()
	if _, hasDefaultLng := res[defaultLng]; !hasDefaultLng {
		res[defaultLng] = o.newLocalizer(bundle, defaultLng)
	}
	return res
}

func (o *options) newLocalizer(bundle *i18n.Bundle, lang string) *i18n.Localizer {
	defaultLang := o.defaultLanguage.String()
	langs := []string{
		lang,
	}
	if lang != defaultLang {
		langs = append(langs, defaultLang)
	}
	localizer := i18n.NewLocalizer(bundle, langs...)
	return localizer
}

// loadFile 加载文件
func (o *options) loadFile(bundle *i18n.Bundle) {
	// 加载配置文件
	for _, tag := range o.acceptLanguage {
		path := filepath.Join(o.languagePath, tag.String()) + "." + o.format
		buf, err := ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}
		if _, err := bundle.ParseMessageFileBytes(buf, path); err != nil {
			panic(err)
		}
	}
}

func defaultOptions() *options {
	return &options{
		languagePath:    "../lang/",
		format:          "toml",
		getLangHandler:  defaultGetLangHandler,
		defaultLanguage: language.English,
		acceptLanguage: []language.Tag{
			language.English,
			language.Chinese,
		},
	}
}

// Format 设置格式化
func Format(format string) Option {
	return func(opts *options) {
		opts.format = format
	}
}

// Path 设置语言文件跟路径
func Path(path string) Option {
	return func(opts *options) {
		opts.languagePath = path
	}
}

// DefaultLang 设置默认
func DefaultLang(lang language.Tag) Option {
	return func(opts *options) {
		opts.defaultLanguage = lang
	}
}

// AddAcceptLanguage 添加accept语言标签
func AddAcceptLanguage(tags ...language.Tag) Option {
	return func(opts *options) {
		opts.acceptLanguage = append(opts.acceptLanguage, tags...)
	}
}
