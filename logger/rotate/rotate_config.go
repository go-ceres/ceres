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
	"time"
)

// RotateConfig 日志分割配置信息
type RotateConfig struct {
	Filename   string        // Filename 			日志文件路径。默认为空，表示关闭，仅输出到终端
	MaxSize    int           // MaxSize 		按照日志文件大小对文件进行滚动切分。默认为0，表示关闭滚动切分特性
	MaxAge     int           // MaxAge			按照切分的文件有效期清理切分文件，当滚动切分特性开启时有效。默认为0，表示不备份，切分则删除
	MaxBackups int           // MaxBackups		文件保存的最大数量，默认值为10
	Interval   time.Duration // Interval		日志轮转的时间，默认当前条件无效
	LocalTime  bool          // LocalTime		LocalTime确定用于格式化中的时间戳的时间,备份文件是计算机的本地时间
	Compress   bool          // Compress		是否使用gzip压缩日志文件，默认不压缩
}

// Build 根据配置文件构建实例
func (rc *RotateConfig) Build() *Rotate {
	rotate := newRotate(rc)
	return rotate
}

func NewDefaultRotateConfig() *RotateConfig {
	return &RotateConfig{
		Filename:   "default.log",
		MaxSize:    500,
		MaxAge:     7,
		MaxBackups: 10,
		Interval:   0,
		LocalTime:  true,
		Compress:   false,
	}
}
