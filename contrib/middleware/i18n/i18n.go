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
	"errors"
	"fmt"
	_ "github.com/go-ceres/ceres/pkg/common/codec/json"
	_ "github.com/go-ceres/ceres/pkg/common/codec/toml"
	_ "github.com/go-ceres/ceres/pkg/common/codec/xml"
	_ "github.com/go-ceres/ceres/pkg/common/codec/yaml"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"sync"
)

var _ I18n = (*I18Impl)(nil)

type I18n interface {
	// GetMessage 获取消息
	getMessage(param interface{}) (string, error)
	// mustGetMessage 忽略错误并获取消息
	mustGetMessage(param interface{}) string
	// setCurrentContext 设置当前上下文
	setCurrentContext(tr context.Context)
}

// I18Impl 实现
type I18Impl struct {
	opts       *options
	bundle     *i18n.Bundle
	ctx        context.Context
	localizers map[string]*i18n.Localizer
	locker     *sync.Mutex
}

// newI18n 创建实例
func newI18n(opt ...Option) I18n {
	opts := defaultOptions()
	for _, fn := range opt {
		fn(opts)
	}
	// 构建bundle
	bundle := opts.buildBundle()
	// 加载语言
	localizers := opts.buildLocalizers(bundle)
	ins := &I18Impl{
		opts:       opts,
		locker:     &sync.Mutex{},
		bundle:     bundle,
		localizers: localizers,
	}
	// 设置
	return ins
}

// setCurrentContext 设置当前上下文
func (i *I18Impl) setCurrentContext(ctx context.Context) {
	i.locker.Lock()
	defer i.locker.Unlock()
	i.ctx = ctx
}

// getMessage 获取消息
func (i *I18Impl) getMessage(param interface{}) (string, error) {
	lang := i.opts.getLangHandler(i.ctx, i.opts.defaultLanguage.String())
	localizer := i.getLocalizerByLang(lang)

	conf, err := i.getLocalizeConfig(param)
	if err != nil {
		return "", err
	}
	return localizer.Localize(conf)
}

// mustGetMessage 忽略错误
func (i *I18Impl) mustGetMessage(param interface{}) string {
	message, _ := i.getMessage(param)
	return message
}

// getLocalizeConfig 获取读取消息配置
func (i *I18Impl) getLocalizeConfig(param interface{}) (*i18n.LocalizeConfig, error) {
	switch paramValue := param.(type) {
	case string:
		localizeConfig := &i18n.LocalizeConfig{
			MessageID: paramValue,
		}
		return localizeConfig, nil
	case *i18n.LocalizeConfig:
		return paramValue, nil
	}
	return nil, errors.New(fmt.Sprintf("un supported localize param: %v", param))
}

// getLocalizerByLang 获取Localizer
func (i *I18Impl) getLocalizerByLang(lang string) *i18n.Localizer {
	localizer, ok := i.localizers[lang]
	if ok {
		return localizer
	}
	return i.localizers[i.opts.defaultLanguage.String()]
}
