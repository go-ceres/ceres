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

package flag

import (
	"errors"
	"fmt"
	"github.com/go-ceres/ceres"
	"github.com/spf13/pflag"
	"os"
)

var showVersion = pflag.BoolP("version", "v", false, "show version")
var showHelp = pflag.BoolP("help", "h", false, "show help")

func init() {
	pflag.ErrHelp = errors.New("help requested")
	pflag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage of %s:\n", ceres.AppName())
		pflag.PrintDefaults()
	}
}

func BoolSliceVar(p *[]bool, name string, value []bool, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.BoolSliceVarP(p, name, short, value, usage)
}

func BoolSlice(name string, value []bool, usage string, shorthand ...string) *[]bool {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.BoolSliceP(name, key, value, usage)
}

func BoolVar(p *bool, name string, value bool, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.BoolVarP(p, name, short, value, usage)
}

func Bool(name string, value bool, usage string, shorthand ...string) *bool {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.BoolP(name, key, value, usage)
}

func Float64SliceVar(p *[]float64, name string, value []float64, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Float64SliceVarP(p, name, short, value, usage)
}

func Float64Slice(name string, value []float64, usage string, shorthand ...string) *[]float64 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Float64SliceP(name, key, value, usage)
}

func Float32SliceVar(p *[]float32, name string, value []float32, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Float32SliceVarP(p, name, short, value, usage)
}

func Float32Slice(name string, value []float32, usage string, shorthand ...string) *[]float32 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Float32SliceP(name, key, value, usage)
}

func Float64Var(p *float64, name string, value float64, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Float64VarP(p, name, short, value, usage)
}

func Float64(name string, value float64, usage string, shorthand ...string) *float64 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Float64P(name, key, value, usage)
}

func Float32Var(p *float32, name string, value float32, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Float32VarP(p, name, short, value, usage)
}

func Float32(name string, value float32, usage string, shorthand ...string) *float32 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Float32P(name, key, value, usage)
}

func UintSliceVar(p *[]uint, name string, value []uint, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.UintSliceVarP(p, name, short, value, usage)
}

func UintSlice(name string, value []uint, usage string, shorthand ...string) *[]uint {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.UintSliceP(name, key, value, usage)
}

func UintVar(p *uint, name string, value uint, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.UintVarP(p, name, short, value, usage)
}

func Uint(name string, value uint, usage string, shorthand ...string) *uint {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.UintP(name, key, value, usage)
}

func Uint64Var(p *uint64, name string, value uint64, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Uint64VarP(p, name, short, value, usage)
}

func Uint64(name string, value uint64, usage string, shorthand ...string) *uint64 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Uint64P(name, key, value, usage)
}

func Uint32Var(p *uint32, name string, value uint32, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Uint32VarP(p, name, short, value, usage)
}

func Uint32(name string, value uint32, usage string, shorthand ...string) *uint32 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Uint32P(name, key, value, usage)
}

func Uint16Var(p *uint16, name string, value uint16, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Uint16VarP(p, name, short, value, usage)
}

func Uint16(name string, value uint16, usage string, shorthand ...string) *uint16 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Uint16P(name, key, value, usage)
}

func Uint8Var(p *uint8, name string, value uint8, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Uint8VarP(p, name, short, value, usage)
}

func Uint8(name string, value uint8, usage string, shorthand ...string) *uint8 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Uint8P(name, key, value, usage)
}

func Int64Var(p *int64, name string, value int64, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Int64VarP(p, name, short, value, usage)
}

func Int64(name string, value int64, usage string, shorthand ...string) *int64 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Int64P(name, key, value, usage)
}

func Int32Var(p *int32, name string, value int32, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Int32VarP(p, name, short, value, usage)
}

func Int32(name string, value int32, usage string, shorthand ...string) *int32 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Int32P(name, key, value, usage)
}

func Int16Var(p *int16, name string, value int16, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Int16VarP(p, name, short, value, usage)
}

func Int16(name string, value int16, usage string, shorthand ...string) *int16 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Int16P(name, key, value, usage)
}

func Int8Var(p *int8, name string, value int8, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.Int8VarP(p, name, short, value, usage)
}

func Int8(name string, value int8, usage string, shorthand ...string) *int8 {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.Int8P(name, key, value, usage)
}

func IntSliceVarP(p *[]int, name string, value []int, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.IntSliceVarP(p, name, short, value, usage)
}

func IntSlice(name string, value []int, usage string, shorthand ...string) *[]int {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.IntSliceP(name, key, value, usage)
}

func IntVar(p *int, name string, value int, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.IntVarP(p, name, short, value, usage)
}

func Int(name string, value int, usage string, shorthand ...string) *int {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.IntP(name, key, value, usage)
}

func StringSliceVar(p *[]string, name string, value []string, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.StringSliceVarP(p, name, short, value, usage)
}

func StringSliceP(name string, value []string, usage string, shorthand ...string) *[]string {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.StringSliceP(name, key, value, usage)
}

func StringVar(p *string, name, value string, usage string, shorthand ...string) {
	short := ""
	if len(shorthand) > 0 {
		short = shorthand[0]
	}
	pflag.StringVarP(p, name, short, value, usage)
}

func String(name, value string, usage string, shorthand ...string) *string {
	key := ""
	if len(shorthand) > 0 {
		key = shorthand[0]
	}
	return pflag.StringP(name, key, value, usage)
}

func Parse() {
	pflag.Parse()
	if *showVersion {
		ceres.ShowVersion()
		os.Exit(0)
	}
	if *showHelp {
		pflag.Usage()
		os.Exit(0)
	}
}
