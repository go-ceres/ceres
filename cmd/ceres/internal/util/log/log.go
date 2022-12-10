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

package log

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"os"
)

type Log struct {
	enable bool // 啰嗦模式
}

// NewLog returns an instance of colorConsole
func NewLog(enable bool) *Log {
	return &Log{
		enable: enable,
	}
}

func (c *Log) Info(format string, a ...interface{}) {
	if !c.enable {
		return
	}
	msg := fmt.Sprintf(format, a...)
	fmt.Println(msg)
}

func (c *Log) Debug(format string, a ...interface{}) {
	if !c.enable {
		return
	}
	msg := fmt.Sprintf(format, a...)
	println(aurora.BrightCyan(msg))
}

func (c *Log) Success(format string, a ...interface{}) {
	if !c.enable {
		return
	}
	msg := fmt.Sprintf(format, a...)
	println(aurora.BrightGreen(msg))
}

func (c *Log) Warning(format string, a ...interface{}) {
	if !c.enable {
		return
	}
	msg := fmt.Sprintf(format, a...)
	println(aurora.BrightYellow(msg))
}

func (c *Log) Error(format string, a ...interface{}) {
	if !c.enable {
		return
	}
	msg := fmt.Sprintf(format, a...)
	println(aurora.BrightRed(msg))
}

func (c *Log) Fatalln(format string, a ...interface{}) {
	if !c.enable {
		return
	}
	c.Error(format, a...)
	os.Exit(1)
}

func (c *Log) MarkDone() {
	if !c.enable {
		return
	}
	c.Success("Done.")
}

func (c *Log) Must(err error) {
	if !c.enable {
		return
	}
	if err != nil {
		c.Fatalln("%+v", err)
	}
}
