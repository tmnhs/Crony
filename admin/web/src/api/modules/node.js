
import request from '@/utils/request'
import { paramsSerializer } from '@/utils/core'

// 获取node节点列表
export const getNodeList = data =>
	request.post('/node/search', data, {
		paramsSerializer(params) {
			return paramsSerializer(params)
		},
	})


// 删除job
export const remove = data => request.post('/node/del', data)