import request from '@/utils/request'

export const getScriptList = data =>
request.post('/script/search', data, {
paramsSerializer(params) {
return paramsSerializer(params)
},
})
// 获取script详情
export const findScript = data => request.get('/script/find', data)

// 修改/添加script
export const updateScript = data => request.post('/script/add', data)

// 删除script
export const removeScript = data => request.post('/script/del', data)