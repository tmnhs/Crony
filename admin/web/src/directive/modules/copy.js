/* 复制一段内容到剪切板板 */

import { Message } from 'element-ui'

const copyNode = elem => {
	const selection = window.getSelection()
	// 如果剪切板中已经有复制了的内容，需要清掉。
	if (selection.rangeCount > 0) {
		selection.removeAllRanges()
	}
	const range = document.createRange()
	range.selectNodeContents(elem)
	selection.addRange(range)
	const result = document.execCommand('Copy')
	// 清除选中的内容,也可以使用 window.getSelection().removeAllRanges()
	range.collapse(false)

	if (result) {
		Message({
			type: 'success',
			message: '复制成功',
		})
	} else {
		Message({
			type: 'error',
			message: '复制失败',
		})
	}
}

const copy = content => {
	if (!content) {
		Message({
			type: 'warning',
			message: '没有要复制的内容',
		})
		return
	}

	if (content.nodeType === 1) {
		copyNode(content)
	} else if (typeof content === 'string') {
		const wrap = document.createElement('div')
		wrap.insertAdjacentText('afterBegin', content)
		document.body.appendChild(wrap)
		copyNode(wrap)
		document.body.removeChild(wrap)
	} else {
		Message({
			type: 'warning',
			message: '没有可以复制的内容',
		})
	}
}

export default {
	inserted(el, binding) {
		el.content = binding.value
		el.addEventListener('click', () => {
			copy(el.content)
		})
	},
	update(el, binding) {
		el.content = binding.value
	},
}

// 获取某个元素中的文本内容
// const text = ele.textContent;
