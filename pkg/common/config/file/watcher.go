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

package file

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/go-ceres/ceres/pkg/common/config"
	"os"
	"path/filepath"
)

type watcher struct {
	f      *file
	ctx    context.Context
	cancel context.CancelFunc
	fw     *fsnotify.Watcher
}

func (w watcher) Next() ([]*config.DataSet, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case event := <-w.fw.Events:
		if event.Op == fsnotify.Rename {
			if _, err := os.Stat(event.Name); err == nil || os.IsExist(err) {
				if err := w.fw.Add(event.Name); err != nil {
					return nil, err
				}
			}
		}
		fi, err := os.Stat(w.f.path)
		if err != nil {
			return nil, err
		}
		path := w.f.path
		if fi.IsDir() {
			path = filepath.Join(w.f.path, filepath.Base(event.Name))
		}
		kv, err := w.f.loadFile(path)
		if err != nil {
			return nil, err
		}
		return []*config.DataSet{kv}, nil
	case err := <-w.fw.Errors:
		return nil, err
	}
}

func (w watcher) Stop() error {
	w.cancel()
	return w.fw.Close()
}

func newWatcher(f *file) (config.Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	if err = w.Add(f.path); err != nil {
		return nil, err
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	return &watcher{
		f:      f,
		fw:     w,
		ctx:    ctx,
		cancel: cancelFunc,
	}, nil
}
