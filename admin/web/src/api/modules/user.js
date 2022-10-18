import request from '@/utils/request'
import { paramsSerializer } from '@/utils/core'

// 获取用户列表
export const getList = data =>
	request.post('/user/search', data, {
		paramsSerializer(params) {
			return paramsSerializer(params)
		},
	})

// 获取用户详情
export const getDetail = data => request.get('/user/find', data)

// 修改用户
export const update = data => request.post('/user/update', data)

// 删除用户
export const remove = data => request.post('/user/del', data)
