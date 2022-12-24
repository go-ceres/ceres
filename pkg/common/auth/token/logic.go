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
	"strconv"
	"time"
)

type Logic struct {
	conf *Options
}

// NewLogic 根据配置创建一个auth处理逻辑
func NewLogic(c ...*Options) *Logic {
	conf := DefaultOptions()
	if len(c) > 0 {
		conf = c[0]
	}
	// 如果没有设置权限管理器
	if conf.permission == nil {
		conf.logger.Fatalf("not set permission manager")
	}
	// 不能没有存储器
	if conf.storage == nil {
		conf.logger.Fatalf("not set storage")
	}
	// 如果没有设置token构建器
	if conf.tokenBuilder == nil {
		conf.tokenBuilder = &defaultTokenBuilder{
			style: conf.TokenStyle,
		}
	}
	// 如果没有设置监听器
	if conf.listener == nil {
		conf.listener = defaultListener{
			logger: conf.logger,
		}
	}

	return &Logic{
		conf: conf,
	}
}

// GetLogicType 获取登录逻辑
func (l *Logic) GetLogicType() string {
	return l.conf.LogicType
}

// ================== 获取token相关 =================

// GetTokenName 获取当前logic的token名
func (l *Logic) GetTokenName() string {
	return l.conf.TokenName
}

// CreateTokenValue 创建token
func (l *Logic) CreateTokenValue(loginId string, device string) string {
	return l.conf.tokenBuilder.Build(loginId, l.conf.LogicType, device)
}

// GetTokenInfo 获取指定token的登录信息
func (l *Logic) GetTokenInfo(tokenValue string) *Info {
	info := new(Info)
	info.Name = l.GetTokenName()
	info.Value = tokenValue
	info.IsLogin = l.IsLogin(tokenValue)
	info.LogicType = l.conf.LogicType
	info.TokenTimeout = l.GetTokenTimeout(tokenValue)
	info.SessionTimeout = l.GetSessionTimeoutByLoginId(l.GetLoginIdNotHandle(tokenValue))
	info.TokenSessionTimeout = l.GetTokenSessionTimeoutByTokenValue(tokenValue)
	info.TokenActivityTimeout = l.GetTokenActivityTimeoutByToken(tokenValue)
	info.LoginDevice = l.GetLoginDevice(tokenValue)
	return info
}

// ================== 登录相关操作 =================

// Login 登录
// 形参:
//
//	loginId – 账号id
//	opts – 设备和超时时间
func (l *Logic) Login(loginId string, opts ...LoginOption) (string, error) {
	// 1.判断用户id
	if loginId == "" {
		return "", ErrorNoToken("账号id不能为空")
	}
	// 2.检查用户是否被禁用
	if l.IsDisable(loginId) {
		return "", ErrorDisableLogin("loginId:%s账号已经被封禁", loginId)
	}
	// 初始化参数
	opt := defaultLoginOptions(l.conf)
	for _, option := range opts {
		option(opt)
	}

	// 3.生成token
	var tokenValue = ""
	// ----如果允许并发登录
	if l.conf.IsConcurrent {
		// ---- 如果配置为共享token，则尝试从Session签名记录中去除token
		if l.conf.IsShare {
			tokenValue = l.GetTokenValueByLoginId(loginId, opt.device)
		}
	} else {
		// 如果不允许并发登录，则将该账号的历史登录标识标记为：被顶下线
		_ = l.Replaced(loginId, opt.device)
	}
	// 如果至此，仍未成功创建tokenValue, 则开始生成一个
	if len(tokenValue) == 0 {
		tokenValue = l.CreateTokenValue(loginId, opt.device)
	}

	// 4. 获取 User-Session , 续期
	ss := l.GetSessionByLoginId(loginId, true)
	ss.UpdateMinTimeout(opt.timeout)
	// 在 User-Session 上记录token签名
	ss.AddTokenSign(tokenValue, opt.device)

	// 4. 持久化其他数据
	// token -> id 映射关系
	l.SaveTokenToIdMapping(tokenValue, loginId, opt.timeout)
	// 写入 [token-last-activity]
	l.setLastActivityToNow(tokenValue)
	// 通知监听器，账号xxx 登录成功
	l.conf.listener.DoLogin(l.conf.LogicType, loginId, *opt)
	return tokenValue, nil
}

// Logout 会话注销，根据账号id 和 设备标识
// 形参:
//
//	loginId – 账号id
//	device – 设备标识 (填""代表所有注销设备)
func (l *Logic) Logout(loginId string, device string) {
	l.clearTokenCommonMethod(loginId, device, func(token string) {
		// 删除Token-Id映射 & 清除Token-Session
		l.DeleteTokenToIdMapping(token)
		l.DeleteTokenSession(token)
		l.conf.listener.DoLogout(l.conf.LogicType, loginId, token)
	}, true)
}

// LogoutByTokenValue 会话注销，根据token
// 形参：
//
//	tokenValue - token值
func (l *Logic) LogoutByTokenValue(tokenValue string) {
	// 1. 清理 token-last-activity
	l.clearLastActivity(tokenValue)
	// 2. 注销 Token-Session
	l.DeleteTokenSession(tokenValue)
	// 3. 获取用户id
	loginId := l.GetLoginIdNotHandle(tokenValue)
	if l.IsValidLoginId(loginId) == false {
		if len(loginId) != 0 {
			l.DeleteTokenToIdMapping(tokenValue)
		}
		return
	}
	// 4. 清理token-id索引
	l.DeleteTokenToIdMapping(tokenValue)
	// 通知监听器
	l.conf.listener.DoLogout(l.conf.LogicType, loginId, tokenValue)
	// 5. 清理User-Session上的token签名 & 尝试注销User-Session
	ss := l.GetSessionByLoginId(loginId, false)
	if len(ss.Id) != 0 {
		ss.RemoveSign(tokenValue)
		ss.LogoutByTokenSignCountToZero()
	}
}

// Replaced 顶人下线，根据账号id 和 设备标识
// 当对方再次访问系统时，会抛出错误，场景值=-4
// 形参:
//
//	loginId – 账号id
//	device – 设备标识 (填null代表顶替所有设备)
func (l *Logic) Replaced(loginId string, device string) error {
	// 如果没有id
	if len(loginId) == 0 {
		return ErrorNoUserId("LoginId 不能为空")
	}
	l.clearTokenCommonMethod(loginId, device, func(token string) {
		// 将此 token 标记为已被顶替
		l.UpdateTokenToIdMapping(token, BeReplaced)
		// 定人下线监听器
		l.conf.listener.DoReplaced(l.conf.LogicType, loginId, token)
	}, false)
	return nil
}

// --- 会话查询

// IsLogin 查询指定token是否登录
// 形参
// tokenValue -- token值
func (l *Logic) IsLogin(tokenValue string) bool {
	return len(l.GetLoginIdDefaultNull(tokenValue)) != 0
}

// CheckLogin 查询指定token是否登录
// 形参
// tokenValue -- token值
func (l *Logic) CheckLogin(tokenValue string) bool {
	loginId, err := l.GetLoginId(tokenValue)
	if err != nil {
		return false
	}
	if len(loginId) == 0 {
		return false
	}
	return true
}

// GetLoginId 根据指定token获取用户，不存在则返回错误信息
// 形参：
// tokenValue - 指定的token
func (l *Logic) GetLoginId(tokenValue string) (string, error) {
	// 查找此token对应loginId, 如果找不到则抛出：无效token
	loginId := l.GetLoginIdByToken(tokenValue)
	switch loginId {
	case "":
		return "", ErrorNotLogin(l.conf.LogicType, InvalidToken, tokenValue)
	case TokenTimeout:
		return "", ErrorNotLogin(l.conf.LogicType, TokenTimeout, tokenValue)
	case BeReplaced:
		return "", ErrorNotLogin(l.conf.LogicType, BeReplaced, tokenValue)
	case KickOut:
		return "", ErrorNotLogin(l.conf.LogicType, KickOut, tokenValue)
	}
	// 检查是否已经 [临时过期]
	err := l.CheckActivityTimeout(tokenValue)
	if err != nil {
		return "", err
	}
	// 如果配置了自动刷新
	if l.conf.AutoRenew {
		l.UpdateLastActivityToNow(tokenValue)
	}
	// 到此可以返回用户了
	return loginId, nil
}

// GetLoginIdByToken 获取指定Token对应的账号id，如果未登录，则返回 ""
// 形参:
//
//	tokenValue – token
func (l *Logic) GetLoginIdByToken(tokenValue string) string {
	if len(tokenValue) == 0 {
		return ""
	}
	return l.GetLoginIdNotHandle(tokenValue)
}

// GetLoginIdNotHandle 获取指定Token对应的账号id (不做任何特殊处理)
// 形参:
// tokenValue – token值
func (l *Logic) GetLoginIdNotHandle(tokenValue string) string {
	return l.conf.storage.Get(l.splicingTokenValueKey(tokenValue))
}

// GetLoginIdDefaultNull 获取当前会话账号id, 如果未登录，则返回"",并返回错误
func (l *Logic) GetLoginIdDefaultNull(tokenValue string) string {
	// 如果连token都是空的，则直接返回
	if len(tokenValue) == 0 {
		return ""
	}
	// 获取登录的用户编号
	loginId := l.GetLoginIdNotHandle(tokenValue)
	if len(loginId) == 0 || AbnormalList[loginId] {
		return ""
	}
	// 如果已经[临时过期]
	if l.GetTokenActivityTimeoutByToken(tokenValue) == NotValueExpire {
		return ""
	}
	return loginId
}

// ---- 其他操作

// IsValidLoginId 判断指定用户id是否有效
// 形参
// loginId - 用户id
func (l *Logic) IsValidLoginId(loginId string) bool {
	return len(loginId) != 0 && !AbnormalList[loginId]
}

// clearTokenCommonMethod 封装 注销、踢人、顶人 三个动作的相同代码（无API含义方法)
// 形参:
// loginId – 账号id
// device – 设备标识
// callBack – 回调函数
// isLogoutSession – 是否注销 User-Session
func (l *Logic) clearTokenCommonMethod(loginId string, device string, callBack func(token string), isLogoutSession bool) {
	// 1.没有获取到session，表示账号并没有登录，则不需要任何操作
	ss := l.GetSessionByLoginId(loginId, false)
	if len(ss.Id) == 0 {
		return
	}
	// 2.循环token签名列表，开始删除相关信息
	for _, s := range ss.SignList {
		if len(device) == 0 || device == s.Device {
			// s1.清理掉[token-last-activity]
			l.clearLastActivity(s.Value)
			// s2.从token列表中移除
			ss.RemoveSign(s.Value)
			// s3.执行回调函数
			callBack(s.Value)
		}
	}
	// 注销user-token
	if isLogoutSession {

	}
}

// UpdateTokenToIdMapping 更改 Token 指向的 账号Id 值
// 形参:
//
//	tokenValue – token值
//	loginId – 新的账号Id值
func (l *Logic) UpdateTokenToIdMapping(tokenValue string, loginId string) {
	l.conf.storage.Update(l.splicingTokenValueKey(tokenValue), loginId)
}

// DeleteTokenToIdMapping 删除 Token-Id 映射
// 形参:
//
//	tokenValue – token值
func (l *Logic) DeleteTokenToIdMapping(tokenValue string) {
	l.conf.storage.Del(l.splicingTokenValueKey(tokenValue))
}

// SaveTokenToIdMapping 存储 Token-Id 映射
// 形参:
// tokenValue – token值
// loginId – 账号id
// timeout – 会话有效期 (单位: 秒)
func (l *Logic) SaveTokenToIdMapping(tokenValue string, loginId string, timeout int64) {
	l.conf.storage.Set(l.splicingTokenValueKey(tokenValue), loginId, timeout)
}

// ================== 账号封禁 =================

// Disable 封禁指定账号到指定时间
// 形参:
// loginId - 账号id
// disableTime - 封禁时长 （-1=永久封禁）
func (l *Logic) Disable(loginId string, disableTime int64) bool {
	return l.conf.storage.Set(l.splicingDisableKey(loginId), disableValue, disableTime)
}

// IsDisable 判断指定账号是否被封禁停用
func (l *Logic) IsDisable(loginId string) bool {
	return l.conf.storage.Get(l.splicingDisableKey(loginId)) != ""
}

// GetDisableTime 获取封禁时间
func (l *Logic) GetDisableTime(loginId string) int64 {
	return l.conf.storage.TTl(l.splicingDisableKey(loginId))
}

// ================== token-session相关 =================

// DeleteTokenSession 删除Token-Session
// 形参:
// tokenValue – token值
func (l *Logic) DeleteTokenSession(tokenValue string) {
	l.conf.storage.Del(l.splicingTokenSessionKey(tokenValue))
}

// ================== user-session相关 =================

// GetSessionBySessionId 根据SessionId获取session对象
func (l *Logic) GetSessionBySessionId(sessionId string, isCreate bool) *Session {
	ss := &Session{storage: l.conf.storage}
	_ = l.conf.storage.GetObject(sessionId, ss)
	// 如果没有获取到session，并且设置了自动创建
	if len(ss.Id) == 0 && isCreate {
		// 创建session
		ss = l.CreateSession(sessionId)
		// 存储到storage
		l.conf.storage.SetObject(ss.Id, ss, l.conf.Timeout)
	}
	return ss
}

// GetSessionByLoginId 获取指定loginId的session
func (l *Logic) GetSessionByLoginId(loginId string, isCreate bool) *Session {
	return l.GetSessionBySessionId(l.splicingSessionKey(loginId), isCreate)
}

// CreateSession 创建一个session
func (l *Logic) CreateSession(sessionId string) *Session {
	return NewSession(sessionId, l.conf.storage)
}

// ================== 【临时有效期】验证相关 =================

// clearLastActivity 清除指定token的 [最后操作时间]
// 形参:
// tokenValue – 指定token
func (l *Logic) clearLastActivity(tokenValue string) {
	// 如果没有传入tokenValue或者配置了不验证最后存活时间
	if len(tokenValue) == 0 || l.conf.ActivityTimeout == NeverExpire {
		return
	}
	// 删除最后操作时间
	l.conf.storage.Del(l.splicingLastActivityTimeKey(tokenValue))
}

// setLastActivityToNow 写入指定token的 [最后操作时间] 为当前时间戳
// 形参:
// tokenValue – 指定token
func (l *Logic) setLastActivityToNow(tokenValue string) {
	if len(tokenValue) == 0 || l.conf.ActivityTimeout == NeverExpire {
		return
	}
	l.conf.storage.Set(l.splicingLastActivityTimeKey(tokenValue), strconv.FormatInt(time.Now().UnixMilli(), 10), l.conf.Timeout)
}

// UpdateLastActivityToNow 续签指定token：(将 [最后操作时间] 更新为当前时间戳)
// 形参:
// tokenValue – 指定token
func (l *Logic) UpdateLastActivityToNow(tokenValue string) {
	// 如果token为空 或者 设置了[永不过期], 则立即返回
	if len(tokenValue) == 0 || l.conf.ActivityTimeout == NeverExpire {
		return
	}
	l.conf.storage.Update(l.splicingLastActivityTimeKey(tokenValue), strconv.FormatInt(time.Now().UnixMilli(), 10))
}

// CheckActivityTimeout 检查指定token 是否已经[临时过期]，如果已经过期则返回错误
// 形参：
// tokenValue - token值
func (l *Logic) CheckActivityTimeout(tokenValue string) error {
	// 如果token为空 或者 设置了[永不过期], 则立即返回
	if len(tokenValue) == 0 || l.conf.ActivityTimeout == NeverExpire {
		return nil
	}
	timeout := l.GetTokenActivityTimeoutByToken(tokenValue)
	// -1 代表此token已经被设置永不过期，无须继续验证
	if timeout == NeverExpire {
		return nil
	}
	// -2 代表已过期，抛出异常
	if timeout == NotValueExpire {
		return ErrorNotLogin(l.conf.LogicType, TokenTimeout, tokenValue)
	}
	return nil
}

// ================== 过期时间相关 =================

// GetTokenTimeout 获取指定token的过期时间
// 形参：
// tokenValue - 指定token
func (l *Logic) GetTokenTimeout(tokenValue string) int64 {
	return l.conf.storage.TTl(l.splicingTokenValueKey(tokenValue))
}

// GetSessionTimeoutByLoginId 获取指定 loginId 的 User-Session 剩余有效时间 (单位: 秒)
// 形参：
// loginId - 指定用户
func (l *Logic) GetSessionTimeoutByLoginId(loginId string) int64 {
	return l.conf.storage.TTl(l.splicingSessionKey(loginId))
}

// GetTokenSessionTimeoutByTokenValue 获取指定 Token-Session 剩余有效时间 (单位: 秒)
// 形参:
// tokenValue – 指定token
func (l *Logic) GetTokenSessionTimeoutByTokenValue(tokenValue string) int64 {
	return l.conf.storage.TTl(l.splicingTokenValueKey(tokenValue))
}

// GetTokenActivityTimeoutByToken 获取指定token的临时有效期
// 形参:
// tokenValue – 指定token
func (l *Logic) GetTokenActivityTimeoutByToken(tokenValue string) int64 {
	// 如果token为空 , 则返回 -2
	if len(tokenValue) == 0 {
		return NotValueExpire
	}
	// 如果设置了永不过期, 则返回 -1
	if l.conf.ActivityTimeout == NeverExpire {
		return NeverExpire
	}
	// 除开特殊情况，就开始查询
	lastActivityTimeKey := l.splicingLastActivityTimeKey(tokenValue)
	lastActivityTimeStr := l.conf.storage.Get(lastActivityTimeKey)
	// 如果查询不到,则返回-2
	if len(lastActivityTimeStr) == 0 {
		return NotValueExpire
	}
	// 计算相差时间
	lastActivityTime, _ := strconv.ParseInt(lastActivityTimeStr, 10, 64) // 存储的时间
	apartSecond := (time.Now().UnixMilli() - lastActivityTime) / 1000    // 当前时间与最后一次活动的时间差值（单位秒）
	timeout := l.conf.ActivityTimeout - apartSecond                      // 配置的最后过期时间差值与实际差值
	// 如果 < 0， 代表已经过期 ，返回-2
	if timeout < 0 {
		return NotValueExpire
	}
	return timeout
}

// ================== 角色验证操作 =================

// GetRoleList 获取：指定账号的角色集合
// 形参:
//
//	loginId – 指定账号id
func (l *Logic) GetRoleList(loginId string) ([]string, error) {
	return l.conf.permission.GetRoleListSlice(loginId, l.conf.LogicType)
}

// HasRole  判断：指定账号是否含有指定角色标识, 返回true或false
// 形参:
//
//	loginId – 账号id
//	role – 角色标识
func (l *Logic) HasRole(loginId string, role string) bool {
	list, err := l.GetRoleList(loginId)
	if err != nil {
		return false
	}
	return contains[string](list, role)
}

// HasRoleAnd 判断：当前账号是否含有指定角色标识 [指定多个，必须全部验证通过]
// 形参:
//
//	loginId - 指定用户
//	roleArray – 角色标识数组
func (l *Logic) HasRoleAnd(loginId string, roleArray ...string) bool {
	list, err := l.GetRoleList(loginId)
	if err != nil {
		return false
	}
	for _, s := range roleArray {
		if !contains(list, s) {
			return false
		}
	}
	return true
}

// CheckRoleOr 校验：当前账号是否含有指定角色标识 [指定多个，只要其一验证通过即可]
// 形参:
//
//	loginId - 指定用户
//	roleArray – 角色标识数组
func (l *Logic) CheckRoleOr(loginId string, roleArray ...string) bool {
	list, err := l.GetRoleList(loginId)
	if err != nil {
		return false
	}
	// 找到任意一个角色就返回
	for _, s := range roleArray {
		if contains(list, s) {
			return true
		}
	}
	return false
}

// ================== 权限验证操作 =================

// GetPermissionList 获取：指定账号的权限码集合
// 形参:
// loginId – 指定账号id
func (l *Logic) GetPermissionList(loginId string) []string {
	permissions, err := l.conf.permission.GetPermissionSlice(loginId, l.conf.LogicType)
	if err != nil {
		return []string{}
	}
	return permissions
}

// HasPermission 判断：当前账号是否含有指定权限, 返回true或false
// 形参:
//
//	loginId - 指定用户
//	permission – 权限码
func (l *Logic) HasPermission(loginId string, permission string) bool {
	slice, err := l.conf.permission.GetPermissionSlice(loginId, l.conf.LogicType)
	if err != nil {
		return false
	}
	return contains(slice, permission)
}

// HasPermissionAnd 判断：指定账号是否含有指定权限, [指定多个，必须全部具有]
// 形参:
// loginId - 指定用户
// permissionArray – 权限码数组
func (l *Logic) HasPermissionAnd(loginId string, permissionArray ...string) bool {
	slice, err := l.conf.permission.GetPermissionSlice(loginId, l.conf.LogicType)
	if err != nil {
		return false
	}
	// 必须包含所有的权限
	for _, s := range permissionArray {
		if !contains(slice, s) {
			return false
		}
	}
	return true
}

// HasPermissionOr 判断：指定账号是否含有指定权限, [指定多个，任意一个有]
// 形参:
// loginId - 指定用户
// permissionArray – 权限码数组
func (l *Logic) HasPermissionOr(loginId string, permissionArray ...string) bool {
	slice, err := l.conf.permission.GetPermissionSlice(loginId, l.conf.LogicType)
	if err != nil {
		return false
	}
	// 有任意一个权限即可
	for _, s := range permissionArray {
		if contains(slice, s) {
			return true
		}
	}
	return false
}

// HasPathPermission 判断：指定账号是否含有指定路由
// 形参:
// loginId - 指定用户
// path – 指定路由路径
func (l *Logic) HasPathPermission(loginId string, path, method string) bool {
	routeInfo, err := l.conf.permission.HasMethodPermission(loginId, l.conf.LogicType, path, method)
	if err != nil {
		return false
	}
	return routeInfo
}

// ================== 根据id反查token相关 =================

// GetLoginDevice 根据指定token值查询登录设备
// 形参：
//
//	tokenValue -- token值
func (l *Logic) GetLoginDevice(tokenValue string) string {
	if len(tokenValue) == 0 {
		return ""
	}
	// 如果是还没有登录
	if !l.IsLogin(tokenValue) {
		return ""
	}
	// 如果session为null的话直接返回 null
	ss := l.GetSessionByLoginId(l.GetLoginIdByToken(tokenValue), false)
	if len(ss.Id) == 0 {
		return ""
	}
	// 遍历解析
	for _, s := range ss.SignList {
		if s.Value == tokenValue {
			return s.Device
		}
	}
	return ""
}

// GetTokenValueByLoginId 获取指定id指定设备端的tokenValue
// 形参
// loginId - 登录用户
// device - 登录设备
func (l *Logic) GetTokenValueByLoginId(loginId string, device string) string {
	// 获取所有的token切片
	tokenSlice := l.GetTokenValueListByLoginId(loginId, device)
	if len(tokenSlice) == 0 {
		return ""
	}
	return tokenSlice[len(tokenSlice)-1]
}

// GetTokenValueListByLoginId 获取指定id指定设备端的token切片
func (l *Logic) GetTokenValueListByLoginId(loginId string, device string) []string {
	ss := l.GetSessionBySessionId(loginId, false)
	// 没有获取到session时
	if len(ss.Id) == 0 {
		return []string{}
	}
	// 遍历
	res := make([]string, 0)
	for _, s := range ss.SignList {
		// 寻找设备相同的token
		if len(device) == 0 || s.Device == device {
			res = append(res, s.Value)
		}
	}
	return res
}

// ================== 返回相应的key =================

// splicingDisableKey 拼接封禁账号key
func (l *Logic) splicingDisableKey(loginId string) string {
	return l.conf.TokenName + ":" + l.conf.LogicType + ":disable:" + loginId
}

// splicingSessionKey 拼接sessionKey
func (l *Logic) splicingSessionKey(loginId string) string {
	return l.GetTokenName() + ":" + l.conf.LogicType + ":session:" + loginId
}

// splicingKeyLastActivityTime
func (l *Logic) splicingLastActivityTimeKey(tokenValue string) string {
	return l.GetTokenName() + ":" + l.conf.LogicType + ":last-activity:" + tokenValue
}

// splicingTokenSessionKey 拼接tokenSessionKey
func (l *Logic) splicingTokenSessionKey(tokenValue string) string {
	return l.GetTokenName() + ":" + l.conf.LogicType + ":token-session:" + tokenValue
}

// splicingTokenValueKey 拼接tokenValueKey
func (l *Logic) splicingTokenValueKey(tokenValue string) string {
	return l.conf.TokenName + ":" + l.conf.LogicType + ":token:" + tokenValue
}

// splicingSwitchKey 拼接切换用户key
func (l *Logic) splicingSwitchKey(tokenValue string) string {
	return l.GetTokenName() + ":" + l.conf.LogicType + ":switch:" + tokenValue
}
