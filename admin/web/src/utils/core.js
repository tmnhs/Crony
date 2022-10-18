/* 核心工具 */

//生成guid
export function guid() {
	let S4 = function () {
		return (((1 + Math.random()) * 0x10000) | 0).toString(16).substring(1)
	}
	return S4() + S4() + '-' + S4() + '-' + S4() + '-' + S4() + '-' + S4() + S4() + S4()
}

// 生成随机数
export function random(min, max) {
	const choice = max - min + 1
	return Math.floor(min + Math.random() * choice)
}

/**
 * 滚动,先慢后快，缓动的效果比easeIn动画明显
 * @param {HTMLDOM} element  要滚动的元素
 * @param {Number} target    目标位置
 * @param {Number} duration  滚动所用的总时间
 * @param {Function} callback  滚动完成之后的回调
 */
export function scroll(element, target, duration = 500, callback = () => {}) {
	const startTime = Date.now()
	const move = () => {
		const passTime = Date.now() - startTime
		const currentPosition = element.scrollTop
		const residueDistance = target - currentPosition
		const step = residueDistance / 10
		element.scrollTop = currentPosition + step
		if (passTime < duration) {
			window.requestAnimationFrame(move)
		} else {
			element.scrollTop = target
			callback()
		}
	}
	move()
}

/**
 * 动画，可改变多个属性
 * @param {HTMLDOM} element  要发生动画的元素
 * @param {Object} properties    要改变的元素属性
 * @param {Number} interval  每次运动的时间间隔
 * @param {Function} callback  动画完成之后的回调
 */
export function animate(element, properties, interval = 20, callback = () => {}) {
	clearInterval(element.timer)
	element.timer = setInterval(() => {
		let flag = true
		for (const property in properties) {
			const current = parseInt(window.getComputedStyle(element)[property])
			const target = properties[property]
			let step = (target - current) / 10
			step = step > 0 ? Math.ceil(step) : Math.floor(step)
			element.style[property] = current + step + 'px'
			if (current != target) {
				flag = false
			}
		}
		if (flag) {
			clearInterval(element.timer)
			callback()
		}
	}, interval)
}

// 用于需要在get请求中传递数组的情况
export function paramsSerializer(params = {}) {
	const paramArr = []
	for (const [key, value] of Object.entries(params)) {
		if (Array.isArray(value)) {
			value.forEach(item => paramArr.push(`${encodeURIComponent(key)}=${encodeURIComponent(item)}`))
		} else {
			paramArr.push(`${encodeURIComponent(key)}=${encodeURIComponent(value)}`)
		}
	}
	return paramArr.join('&')
}

/**
 * 获取url中的查询字符串参数
 * @param {String} url  url字符串
 */
export function getURLParams(url) {
	const search = url.split('?')[1]
	if (!search) {
		return {}
	}
	return JSON.parse(
		'{"' + decodeURIComponent(search).replace(/"/g, '\\"').replace(/&/g, '","').replace(/=/g, '":"') + '"}'
	)
}

// 深克隆
export function deepClone(source) {
	if (typeof source !== 'object' || source === null) {
		return source
	}
	const target = Array.isArray(source) ? [] : {}
	for (const [key, value] of Object.entries(source)) {
		target[key] = deepClone(value)
	}
	return target
}

// 获取元素相对于浏览器窗口边缘的的距离
export function getOffset(elem) {
	function getLeft(o) {
		if (o == null) {
			return 0
		} else {
			return o.offsetLeft + getLeft(o.offsetParent) + (o.offsetParent ? o.offsetParent.clientLeft : 0)
		}
	}

	function getTop(o) {
		if (o == null) {
			return 0
		} else {
			return o.offsetTop + getTop(o.offsetParent) + (o.offsetParent ? o.offsetParent.clientTop : 0)
		}
	}
	return { left: getLeft(elem), top: getTop(elem) }
}

// 节流
export function throttle(fn, interval = 100) {
	let timer = null
	return function () {
		const context = this
		const args = arguments
		if (!timer) {
			timer = setTimeout(() => {
				timer = null
				fn.apply(context, args)
			}, interval)
		}
	}
}

// 防抖
export function debounce(fn, interval = 100) {
	let timer = null
	return function () {
		const context = this
		const args = arguments
		if (timer) {
			clearTimeout(timer)
		}
		timer = setTimeout(() => {
			fn.apply(context, args)
		}, interval)
	}
}

// 判断数据类型
export const getType = value => (value ? value.constructor.name.toLowerCase() : value)

// 加载第三方脚本
export function loadScript(src, callback = (err, res) => {}) {
	const existScript = document.getElementById(src)
	if (existScript) {
		callback(null, existScript)
	} else {
		const script = document.createElement('script')
		script.src = src
		script.id = src
		document.body.appendChild(script)
		script.onload = function () {
			callback(null, script)
		}
		script.onerror = function () {
			callback(new Error(`“${src}”加载失败`), script)
		}
	}
}

/**
 * 获取树的所有节点的某个属性值
 */
export const getTreeNodeValue = (tree, filed) => {
	return tree
		.map(node => {
			const result = []
			node[filed] && result.push(node[filed])
			if (node.children) {
				result.push(...getTreeNodeValue(node.children, filed))
			}
			return result
		})
		.flat()
}
