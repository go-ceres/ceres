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

package rotate

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"syscall"
	"time"
)

// Rotate实例，实际上是一个io.WriteCloser，去实现writer和close方法
var (
	_ io.WriteCloser = (*Rotate)(nil)
	// 获取当前时间的函数
	currentTime = time.Now
	// 获取文件状态信息的函数
	osStat = os.Stat
)

type (
	Rotate struct {
		config    *RotateConfig
		size      int64
		ctime     time.Time
		file      *os.File
		mu        sync.Mutex
		millCh    chan bool
		startMill sync.Once
	}
	logInfo struct {
		timestamp time.Time
		os.FileInfo
	}
)

const (
	megaByte         = 1024 * 1024
	defaultMaxSize   = 500
	compressSuffix   = ".gz"
	backupTimeFormat = "2006-01-02T15-04-05.000"
)

// 创建一个日志切割实例
func newRotate(config *RotateConfig) *Rotate {
	rotate := Rotate{
		config: config,
	}
	return &rotate
}

func (r *Rotate) Name() string {
	return "file"
}

// Write 写入日志方法，上层调用
func (r *Rotate) Write(p []byte) (n int, err error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	writeLen := int64(len(p))
	if writeLen > r.max() {
		return 0, fmt.Errorf(
			"write length %d exceeds maximum file size %d", writeLen, r.max(),
		)
	}

	if r.file == nil {
		if err = r.openExistingOrNew(len(p)); err != nil {
			return 0, err
		}
	}

	if r.size+writeLen > r.max() {
		if err := r.rotate(); err != nil {
			return 0, err
		}
	}

	if r.config.Interval > 0 {
		cutoff := currentTime().Add(-1 * r.config.Interval)
		if r.ctime.Before(cutoff) {
			if err := r.rotate(); err != nil {
				return 0, err
			}
		}
	}

	n, err = r.file.Write(p)
	r.size += int64(n)
	return n, err
}

// max 计算每个文件的最大存储量
func (r *Rotate) max() int64 {
	if r.config.MaxSize == 0 {
		return int64(defaultMaxSize * megaByte)
	}
	return int64(r.config.MaxSize) * int64(megaByte)
}

// openExistingOrNew 如果存在文件则打开并
func (r *Rotate) openExistingOrNew(writeLen int) error {
	r.mill()
	filename := r.filename()
	info, err := osStat(filename)
	if os.IsNotExist(err) {
		return r.openNew()
	}
	if err != nil {
		return fmt.Errorf("error getting log file info: %s", err)
	}

	if info.Size()+int64(writeLen) >= r.max() {
		return r.rotate()
	}

	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		// if we fail to open the old log file for some reason, just ignore
		// it and open a new log file.
		return r.openNew()
	}
	r.file = file
	r.size = info.Size()
	if ct, err := ctime(file); err == nil {
		r.ctime = ct
	}
	return nil
}

// dir 返回当前文件名的目录。
func (r *Rotate) dir() string {
	return filepath.Dir(r.filename())
}

// prefixAndExt 获取日志文件的后缀，和name
func (r *Rotate) prefixAndExt() (prefix, ext string) {
	filename := filepath.Base(r.filename())
	ext = filepath.Ext(filename)
	prefix = filename[:len(filename)-len(ext)]
	return prefix, ext
}

// oldLogFiles 获取和日志文件同级目录的旧日志文件
// 目录作为当前日志文件，按ModTime排序
func (r *Rotate) oldLogFiles() ([]logInfo, error) {
	files, err := ioutil.ReadDir(r.dir())
	if err != nil {
		return nil, fmt.Errorf("can't read log file directory: %s", err)
	}
	var logFiles []logInfo

	prefix, ext := r.prefixAndExt()

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if t, err := r.timeFromName(f.Name(), prefix, ext); err == nil {
			logFiles = append(logFiles, logInfo{t, f})
			continue
		}
		// error parsing means that the suffix at the end was not generated
		// by rotate, and therefore it's not a backup file.
	}

	sort.Sort(byFormatTime(logFiles))

	return logFiles, nil
}

// millRunOnce 执行过时日志文件的压缩和删除。
// 最多只会保留配置文件配置的MaxBackups数量的日志文件
// 最多只会保留配置文件MaxAge时间的日志
// 最后对符合规范的日志进行处理
func (r *Rotate) millRunOnce() error {
	if r.config.MaxBackups == 0 && r.config.MaxAge == 0 && !r.config.Compress {
		return nil
	}
	files, err := r.oldLogFiles()
	if err != nil {
		return err
	}

	var compress, remove []logInfo

	if r.config.MaxBackups > 0 && r.config.MaxBackups < len(files) {
		preserved := make(map[string]bool)
		var remaining []logInfo
		for _, f := range files {
			// Only count the uncompressed log file or the
			// compressed log file, not both.
			fn := f.Name()
			if strings.HasSuffix(fn, compressSuffix) {
				fn = fn[:len(fn)-len(compressSuffix)]
			}
			preserved[fn] = true

			if len(preserved) > r.config.MaxBackups {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}
	if r.config.MaxAge > 0 {
		diff := time.Duration(int64(24*time.Hour) * int64(r.config.MaxAge))
		cutoff := currentTime().Add(-1 * diff)

		var remaining []logInfo
		for _, f := range files {
			if f.timestamp.Before(cutoff) {
				remove = append(remove, f)
			} else {
				remaining = append(remaining, f)
			}
		}
		files = remaining
	}

	if r.config.Compress {
		for _, f := range files {
			if !strings.HasSuffix(f.Name(), compressSuffix) {
				compress = append(compress, f)
			}
		}
	}

	for _, f := range remove {
		errRemove := os.Remove(filepath.Join(r.dir(), f.Name()))
		if err == nil && errRemove != nil {
			err = errRemove
		}
	}
	for _, f := range compress {
		fn := filepath.Join(r.dir(), f.Name())
		errCompress := compressLogFile(fn, fn+compressSuffix)
		if err == nil && errCompress != nil {
			err = errCompress
		}
	}

	return err
}

// 运行在一个goroutine中来管理轮换后的压缩和移除旧日志文件
func (r *Rotate) millRun() {
	for range r.millCh {
		// what am I going to do, log this?
		_ = r.millRunOnce()
	}
}

// 异步进行日志切割，压缩
func (r *Rotate) mill() {
	r.startMill.Do(func() {
		r.millCh = make(chan bool, 1)
		go r.millRun()
	})
	select {
	case r.millCh <- true:
	default:
	}
}

// filename 获取文件名
func (r *Rotate) filename() string {
	if r.config.Filename != "" {
		return r.config.Filename
	}
	name := filepath.Base(os.Args[0]) + "-rotate.log"
	return filepath.Join(os.TempDir(), name)
}

// timeFromName 通过实践格式化格式获取文件名
func (r *Rotate) timeFromName(filename, prefix, ext string) (time.Time, error) {
	if filename == prefix+ext {
		return time.Time{}, fmt.Errorf("not old file")
	}
	if !strings.HasPrefix(filename, prefix+ext) {
		return time.Time{}, fmt.Errorf("mismatched prefix")
	}
	var ts string
	if !strings.HasSuffix(filename, compressSuffix) {
		ts = filename[len(prefix)+len(ext)+1:]
	} else {
		ts = filename[len(prefix)+len(ext)+1 : len(filename)-len(compressSuffix)]
	}
	return time.Parse(backupTimeFormat, ts)
}

// Close 实现io.WriterClose的close方法
func (r *Rotate) Close() (err error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.close()
}

// close 具体实现关闭方法，即关闭打开的文件
func (r *Rotate) close() error {
	if r.file == nil {
		return nil
	}
	err := r.file.Close()
	r.file = nil
	return err
}

// openNew 打开一个新的文件
func (r *Rotate) openNew() error {
	err := os.MkdirAll(r.dir(), 0755)
	if err != nil {
		return fmt.Errorf("can't make directories for new logfile: %s", err)
	}

	name := r.filename()
	mode := os.FileMode(0644)
	info, err := osStat(name)
	if err == nil {
		// Copy the mode off the old logfile.
		mode = info.Mode()
		// move the existing file
		newname := backupName(name, r.config.LocalTime)
		if err := os.Rename(name, newname); err != nil {
			return fmt.Errorf("can't rename log file: %s", err)
		}

		// this is a no-op anywhere but linux
		if err := chown(name, info); err != nil {
			return err
		}
	}

	// we use truncate here because this should only get called when we've moved
	// the file ourselves. if someone else creates the file in the meantime,
	// just wipe out the contents.
	f, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return fmt.Errorf("can't open new logfile: %s", err)
	}
	r.file = f
	r.size = 0
	r.ctime = currentTime()
	return nil
}

// Rotate 启动日志切割，
func (r *Rotate) Rotate() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.rotate()
}

// rotate 启动日志切割，
func (r *Rotate) rotate() error {
	if err := r.close(); err != nil {
		return err
	}
	if err := r.openNew(); err != nil {
		return err
	}
	r.mill()
	return nil
}

// ctime 获取日志文件的创建时间
func ctime(file *os.File) (time.Time, error) {
	fi, err := file.Stat()
	if err != nil {
		return time.Now(), err
	}
	stat := fi.Sys().(*syscall.Stat_t)
	return time.Unix(int64(stat.Ctimespec.Sec), int64(stat.Ctimespec.Nsec)), nil
}

// backupName 获取日志切割名称，即是备份名
func backupName(name string, local bool) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	t := currentTime()
	if !local {
		t = t.UTC()
	}

	timestamp := t.Format(backupTimeFormat)
	return filepath.Join(dir, fmt.Sprintf("%s%s.%s", prefix, ext, timestamp))
}

//compressLogFile 给指定文件压缩或者删除
func compressLogFile(src, dst string) (err error) {
	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("failed to open log file: %v", err)
	}
	defer func() {
		_ = f.Close()
	}()

	fi, err := osStat(src)
	if err != nil {
		return fmt.Errorf("failed to stat log file: %v", err)
	}

	if err := chown(dst, fi); err != nil {
		return fmt.Errorf("failed to chown compressed log file: %v", err)
	}

	// If this file already exists, we presume it was created by
	// a previous attempt to compress the log file.
	gzf, err := os.OpenFile(dst, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, fi.Mode())
	if err != nil {
		return fmt.Errorf("failed to open compressed log file: %v", err)
	}
	defer func() {
		_ = gzf.Close()
	}()

	gz := gzip.NewWriter(gzf)

	defer func() {
		if err != nil {
			_ = os.Remove(dst)
			err = fmt.Errorf("failed to compress log file: %v", err)
		}
	}()

	if _, err := io.Copy(gz, f); err != nil {
		return err
	}
	if err := gz.Close(); err != nil {
		return err
	}
	if err := gzf.Close(); err != nil {
		return err
	}

	if err := f.Close(); err != nil {
		return err
	}

	return os.Remove(src)
}

// byFormatTime 按名称中格式化的最新时间排序。
type byFormatTime []logInfo

// Less ...
func (b byFormatTime) Less(i, j int) bool {
	return b[i].timestamp.After(b[j].timestamp)
}

// Swap ...
func (b byFormatTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// Len ...
func (b byFormatTime) Len() int {
	return len(b)
}
