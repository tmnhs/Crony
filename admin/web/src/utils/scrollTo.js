// 控制页面窗口或DOM元素滚动

// 缓动公式
const easeIn = (t, b, c, d, s = 4) => {
	return c * (t /= d) * Math.pow(t, s) + b
}

// 获取元素的滚动位置。
const getScrolltPosition = element => {
	if (element === window) {
		return window.pageYOffset || document.body.scrollTop || document.documentElement.scrollTop
	} else {
		return element.scrollTop
	}
}

// 设置要滚动到的位置
const setTargetPosition = (element, target) => {
	if (element === window) {
		document.body.scrollTop = target
		document.documentElement.scrollTop = target
	} else {
		element.scrollTop = target
	}
}

/**
 * 滚动
 * @param {Element | Window} element 要滚动的元素
 * @param {Number} target 要滚动到的最终位置
 * @param {Number} duration    滚动持续时间
 * @param { Function} callback  滚动完成后的回调函数
 */
const scrollTo = (element, target, duration = 500, callback = () => {}) => {
	const startTime = Date.now()
	const startPosition = getScrolltPosition(element)
	const distance = target - startPosition
	const animateScroll = () => {
		const passTime = Date.now() - startTime
		const nextPosition = easeIn(passTime, startPosition, distance, duration)
		setTargetPosition(element, nextPosition)
		if (passTime < duration) {
			window.requestAnimationFrame(animateScroll)
		} else {
			//虽然经过 duration的时间之后元素已经非常接近要滚动到的最终位置，这里再精确设置一下。
			setTargetPosition(element, target)
			callback()
		}
	}
	animateScroll()
}

export default scrollTo
