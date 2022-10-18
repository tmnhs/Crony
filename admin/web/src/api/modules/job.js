
import request from '@/utils/request'
import { paramsSerializer } from '@/utils/core'

// 获取job日志列表
export const getLogList = data =>
	request.post('/job/log', data, {
		paramsSerializer(params) {
			return paramsSerializer(params)
		},
	})
export const getJobList = data =>
request.post('/job/search', data, {
    paramsSerializer(params) {
        return paramsSerializer(params)
    },
})
// 获取job详情
export const findJob = data => request.get('/job/find', data)

// 修改/添加job
export const updateJob = data => request.post('/job/add', data)

// 删除job
export const removeJob = data => request.post('/job/del', data)

// 立即执行job
export const once = data => request.post('/job/once', data)

export const kill = data => request.post('/job/kill', data)

