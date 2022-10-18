/**
 * 计算模糊时间
 */
// type默认为pass表示计算经过了多少时间
export const getSimpleTime = (date, type = 'pass', simple = '') => {
	const second = 1000
	const minute = second * 60
	const hour = minute * 60
	const day = hour * 24
	const month = day * 30
	const year = month * 12
	const now = new Date()
	const target = new Date(date)
	let interval = 0
	let simpleTime = ''
	if (type === 'pass') {
		interval = now - target
	} else {
		interval = target - now
	}
	if (interval < 0) {
		if (type === 'pass') {
			simpleTime = simple || '未开始'
		} else {
			simpleTime = simple || '已结束'
		}
	} else if (0 < interval && interval <= minute) {
		simpleTime = Math.floor(interval / second) + '秒'
	} else if (minute < interval && interval <= hour) {
		simpleTime = Math.floor(interval / minute) + '分钟'
	} else if (hour < interval && interval <= day) {
		simpleTime = Math.floor(interval / hour) + '小时'
	} else if (day < interval && interval <= month) {
		simpleTime = Math.floor(interval / day) + '天'
	} else if (month < interval && interval <= year) {
		simpleTime = Math.floor(interval / month) + '月'
	} else if (year < interval) {
		simpleTime = Math.floor(interval / year) + '年'
	} else {
		simpleTime = `${target.getFullYear()}-${target.getMonth()}-${target.getDate()}`
	}
	if (interval > 0 && type === 'pass') {
		simpleTime = simpleTime + '前'
	}
	return simpleTime
}

export function formatDate(date, fmt) {
	if (fmt==""){
		fmt='yyyy-MM-dd hh:mm'
	}
	if (/(y+)/.test(fmt)) {
	  fmt = fmt.replace(RegExp.$1, (date.getFullYear() + '').substr(4 - RegExp.$1.length))
	}
	let o = {
	  'M+': date.getMonth() + 1,
	  'd+': date.getDate(),
	  'h+': date.getHours(),
	  'm+': date.getMinutes(),
	  's+': date.getSeconds()
	}
	for (let k in o) {
	  if (new RegExp(`(${k})`).test(fmt)) {
		let str = o[k] + ''
		fmt = fmt.replace(RegExp.$1, (RegExp.$1.length === 1) ? str : padLeftZero(str))
	  }
	}
	return fmt
  }
  
  function padLeftZero(str) {
	return ('00' + str).substr(str.length)
  }