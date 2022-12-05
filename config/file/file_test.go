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
	"github.com/go-ceres/ceres/config"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

const (
	_json = `{
    "head":{
        "bar":{
            "string":"ceres",
            "number":123456789123456112,
            "array":[
                "a",
                "b",
                "c"
            ],
            "bool":true,
            "float":12345678.11,
            "time":123456789
        }
    },
    "foot":{
        "mv":{
            "string":"ceres",
            "number":123456789123456112,
            "array":[
                "a",
                "b",
                "c"
            ],
            "bool":true,
            "float":12345678.123
        }
    }
}
`
	_yaml = `head:
  bar:
    string: ceres
    number: 123456789123456112
    array:
      - a
      - b
      - c
    bool: true
    float: 12345678
    time: 123456789
foot:
  mv:
    string: ceres
    number: 123456789123456112
    array:
      - a
      - b
      - c
    bool: true
    float: 12345678
`
	_toml = `[foot.mv]
array = ["a","b","c"]
bool = true
float = 12345678
number = 123456789123456112
string = "ceres"
[head.bar]
array = ["a","b","c"]
bool = true
float = 12345678
number = 123456789123456112
string = "ceres"
time = 123456789
`
)

func TestFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "config")
	jsonFile := filepath.Join(path, "json.json")
	jsonData := []byte(_json)
	defer os.Remove(jsonFile)
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile(jsonFile, jsonData, 0o666); err != nil {
		t.Error(err)
	}
	testSource(t, jsonFile, jsonData)
	testSource(t, path, jsonData)
	testWatchFile(t, jsonFile)
	testWatchDir(t, path, jsonFile)
}

func TestConfig(t *testing.T) {

	path := filepath.Join(t.TempDir(), "config")
	defer os.Remove(path)
	if err := os.MkdirAll(path, 0o700); err != nil {
		t.Error(err)
	}
	if err := os.WriteFile(filepath.Join(path, "test_json.json"), []byte(_json), 0o666); err != nil {
		t.Error(err)
	}
	c, err := config.New(config.WithSource(
		NewSource(path),
	))
	if err != nil {
		t.Error("new config error", err)
	}
	testScan(t, c)
	testConfig(t, c)
}

func testSource(t *testing.T, path string, data []byte) {
	t.Log(path)
	s := NewSource(path)
	kvs, err := s.Load()
	if err != nil {
		t.Error(err)
	}
	if string(kvs[0].Data) != string(data) {
		t.Errorf("no expected: %s, but got: %s", kvs[0].Data, data)
	}
}

func testWatchFile(t *testing.T, path string) {
	t.Log(path)

	s := NewSource(path)
	watch, err := s.Watch()
	if err != nil {
		t.Error(err)
	}

	f, err := os.OpenFile(path, os.O_RDWR, 0)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	_, err = f.WriteString(_json)
	if err != nil {
		t.Error(err)
	}
	kvs, err := watch.Next()
	if err != nil {
		t.Errorf("watch.Next() error(%v)", err)
	}
	if !reflect.DeepEqual(string(kvs[0].Data), _json) {
		t.Errorf("string(kvs[0].Value(%v) is  not equal to _testJSONUpdate(%v)", kvs[0].Data, _json)
	}

	newFilepath := filepath.Join(filepath.Dir(path), "test1.json")
	if err = os.Rename(path, newFilepath); err != nil {
		t.Error(err)
	}
	kvs, err = watch.Next()
	if err == nil {
		t.Errorf("watch.Next() error(%v)", err)
	}
	if kvs != nil {
		t.Errorf("watch.Next() error(%v)", err)
	}

	err = watch.Stop()
	if err != nil {
		t.Errorf("watch.Stop() error(%v)", err)
	}

	if err := os.Rename(newFilepath, path); err != nil {
		t.Error(err)
	}
}

func testWatchDir(t *testing.T, path, file string) {
	t.Log(path)
	t.Log(file)

	s := NewSource(path)
	watch, err := s.Watch()
	if err != nil {
		t.Error(err)
	}

	f, err := os.OpenFile(file, os.O_RDWR, 0)
	if err != nil {
		t.Error(err)
	}
	defer f.Close()
	_, err = f.WriteString(_json)
	if err != nil {
		t.Error(err)
	}

	kvs, err := watch.Next()
	if err != nil {
		t.Errorf("watch.Next() error(%v)", err)
	}
	if !reflect.DeepEqual(string(kvs[0].Data), _json) {
		t.Errorf("string(kvs[0].Value(%s) is  not equal to _testJSONUpdate(%v)", kvs[0].Data, _json)
	}
}

func testScan(t *testing.T, c config.Config) {
	expected := map[string]interface{}{
		"head.bar.number": int64(123456789123456112),
		"head.bar.float":  12345678.11,
		"head.bar.string": "ceres",
		"head.bar.time":   time.Duration(123456789),
	}
	for key, value := range expected {
		switch value.(type) {
		case int64:
			if v, err := c.Get(key).Int(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case float64:
			if v, err := c.Get(key).Float(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case string:
			if v, err := c.Get(key).String(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		case time.Duration:
			if v, err := c.Get(key).Duration(); err != nil {
				t.Error(key, value, err)
			} else if v != value {
				t.Errorf("no expect key: %s value: %v, but got: %v", key, value, v)
			}
		}
	}
	// scan
	var settings struct {
		Number int64         `json:"number"`
		Float  float64       `json:"float"`
		String string        `json:"string"`
		Time   time.Duration `json:"time"`
	}
	if err := c.Get("head.bar").Scan(&settings); err != nil {
		t.Error(err)
	}
	if v := expected["head.bar.number"]; settings.Number != v {
		t.Errorf("no expect number value: %v, but got: %v", settings.Number, v)
	}
	if v := expected["head.bar.float"]; settings.Float != v {
		t.Errorf("no expect float value: %v, but got: %v", settings.Float, v)
	}
	if v := expected["head.bar.string"]; settings.String != v {
		t.Errorf("no expect string value: %v, but got: %v", settings.String, v)
	}
	if v := expected["head.bar.time"]; settings.Time != v {
		t.Errorf("no expect time value: %v, but got: %v", settings.Time, v)
	}
}

func testConfig(t *testing.T, c config.Config) {
	type TestJSON struct {
		Head struct {
			Bar struct {
				Number int64         `json:"number"`
				Float  float64       `json:"float"`
				String string        `json:"string"`
				Time   time.Duration `json:"time"`
			} `json:"bar"`
			Server struct {
				Addr string `json:"addr"`
				Port int    `json:"port"`
			} `json:"server"`
		} `json:"head"`
		Foot struct {
			Name string `json:"name"`
			Age  int    `json:"age"`
		} `json:"foot"`
	}
	var conf TestJSON
	err := c.Scan(&conf)
	if err != nil {
		t.Error(err)
	}
	t.Log(conf)
}
