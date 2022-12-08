//    Copyright 2022. ceres
//    Author https://github.com/go-ceres/ceres
//
//    Licensed under the Apache License, CeresVersion 2.0 (the "License");
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

package version

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/google/uuid"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// CeresVersion ceres的版本
const CeresVersion = "v0.0.4-rc1"

// 变化量
var (
	appId     string // 当前应用唯一标识
	startTime string // 项目开始时间
	goVersion string // go版本号
)

// 构建信息
var (
	appName     string // 应用名称
	hostName    string // 主机名
	appVersion  string // 应用版本
	buildTime   string // 构建时间
	buildUser   string // 构建的用户
	buildStatus string // 应用构建状态
	buildHost   string // 构建的主机
)

// 环境变量设置
var (
	region string
	zone   string
)

// init 初始化
func init() {
	appId = uuid.New().String()
	if appName == "" {
		appName = os.Getenv("APP_NAME")
		if appName == "" {
			appName = filepath.Base(os.Args[0])
		}
	}
	name, err := os.Hostname()
	if err != nil {
		name = "unknown"
	}
	hostName = name
	startTime = time.Now().Format("2006-01-02 15:04:05")
	buildTime = strings.Replace(buildTime, "--", " ", 1)
	goVersion = runtime.Version()
}

// AppId 当前应用唯一标识
func AppId() string {
	return appId
}

// AppName 应用名称
func AppName() string {
	return appName
}

// AppVersion 应用版本
func AppVersion() string {
	return appVersion
}

// BuildTime 构建时间
func BuildTime() string {
	return buildTime
}

// HostName 主机名
func HostName() string {
	return hostName
}

// GoVersion go的运行版本
func GoVersion() string {
	return goVersion
}

// StartTime 项目启动时间
func StartTime() string {
	return startTime
}

// SetAppRegion 设置部属区域
func SetAppRegion(region string) {
	region = region
}

// AppRegion 部属地域
func AppRegion() string {
	return region
}

func SetAppZone(zone string) {
	zone = zone
}

// AppZone 应用部属的分区
func AppZone() string {
	return zone
}

func ShowVersion() {
	fmt.Printf("%-9s]> %-30s => %s\n", "ceres", color.RedString("name"), color.BlueString(appName))
	fmt.Printf("%-9s]> %-30s => %s\n", "ceres", color.RedString("version"), color.BlueString(appVersion))
	fmt.Printf("%-9s]> %-30s => %s\n", "ceres", color.RedString("hostname"), color.BlueString(hostName))
	fmt.Printf("%-9s]> %-30s => %s\n", "ceres", color.RedString("ceresVersion"), color.BlueString(CeresVersion))
	fmt.Printf("%-9s]> %-30s => %s\n", "ceres", color.RedString("goVersion"), color.BlueString(goVersion))
	fmt.Printf("%-9s]> %-30s => %s\n", "ceres", color.RedString("buildUser"), color.BlueString(buildUser))
	fmt.Printf("%-9s]> %-30s => %s\n", "ceres", color.RedString("buildHost"), color.BlueString(buildHost))
	fmt.Printf("%-9s]> %-30s => %s\n", "ceres", color.RedString("buildTime"), color.BlueString(buildTime))
	fmt.Printf("%-9s]> %-30s => %s\n", "ceres", color.RedString("buildStatus"), color.BlueString(buildStatus))
}
