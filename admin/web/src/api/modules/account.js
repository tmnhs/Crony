import request from '@/utils/request'

// 登录
export const login = data => request.post('/login', data)

// 退出登录
export const logout = data => request.post('/logout', data)

// 获取用户信息
export const getUserInfo = data => request.get('/user/find', data)

// 注册
export const register = data => request.post('/register', data)

//  修改密码
export const changePassword = data => request.post('/user/change_pw', data)
