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

import "github.com/casbin/casbin/v2"

var _ Permission = (*CasbinPermission)(nil)

// Permission 权限管理接口
type Permission interface {
	// GetPermissionSlice 获取指定账号指定设备的权限列表
	GetPermissionSlice(loginId string, logicType string) ([]string, error)
	// GetRoleListSlice 获取指定
	GetRoleListSlice(loginId string, logicType string) ([]string, error)
	// HasMethodPermission 检查权限
	HasMethodPermission(loginId, logicType, path, method string) (bool, error)
}

// CasbinPermission 创建使用casbin作为权限管理
type CasbinPermission struct {
	enforcer casbin.IEnforcer
}

func NewCasbinPermission(enforcer casbin.IEnforcer) *CasbinPermission {
	return &CasbinPermission{
		enforcer: enforcer,
	}
}

func (c *CasbinPermission) HasMethodPermission(loginId, logicType, path, method string) (bool, error) {
	return c.enforcer.Enforce(loginId, logicType, path, method)

}

func (c *CasbinPermission) GetPermissionSlice(loginId string, logicType string) ([]string, error) {
	ret := c.enforcer.GetPermissionsForUserInDomain(loginId, logicType)
	res := make([]string, 0, len(ret))
	for _, strings := range ret {
		res = append(res, strings[2])
	}
	return []string{}, nil
}

func (c *CasbinPermission) GetRoleListSlice(loginId string, logicType string) ([]string, error) {
	return c.enforcer.GetRolesForUserInDomain(loginId, logicType), nil
}
