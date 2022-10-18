const setTitle = (el, text = '-') => {
	el.innerText = text
	const isOverFlow = el.scrollWidth > el.clientWidth
	if (isOverFlow) {
		el.setAttribute('title', text)
	} else {
		el.removeAttribute('title')
	}
}

export default {
	bind(el) {
		el.classList.add('text-ellipsis-single')
	},
	inserted(el, binding) {
		setTitle(el, binding.value)
	},
	update(el, binding) {
		setTitle(el, binding.value)
	},
}
