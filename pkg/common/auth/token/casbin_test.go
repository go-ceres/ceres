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

package token

import (
	"github.com/casbin/casbin/v2"
	"testing"
)

func TestCasbin(t *testing.T) {
	enforcer, err := casbin.NewEnforcer("model.conf", "permissions.csv")
	if err != nil {
		t.Error(err)
	}
	enforcer.AddFunction("act_check", func(arguments ...interface{}) (interface{}, error) {
		if len(arguments) == 2 {
			var param string
			var rule string
			switch arguments[0].(type) {
			case string:
				param = arguments[0].(string)
			default:
				return false, nil
			}
			switch arguments[1].(type) {
			case string:
				rule = arguments[1].(string)
			default:
				return false, nil
			}
			if rule == "*" {
				return true, nil
			} else {
				return param == rule, nil
			}
		}
		return false, nil
	})
	res := enforcer.GetPermissionsForUserInDomain("alice", "domain1")
	print(res)
	enforce, err := enforcer.Enforce("alice", "domain1", "/data/1", "read")
	if err != nil {
		t.Error(err)
	}
	type UserInfo struct {
		RoleId int64
	}
	user := &UserInfo{
		RoleId: 1,
	}
	// 获取session
	logic := Logic{}
	session := logic.GetSessionByLoginId("1", true)
	session.DataMap["userInfo"] = user
	session.Update()
	print(enforce)

}
