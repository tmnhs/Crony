import request from '@/utils/request'



// 获取首页数据
export const getTodayStatics = data =>request.get('/statis/today', data)

export const getWeekStatics = data =>request.get('/statis/week', data)


export const getSystemInfo = data =>request.get('/statis/system', data)