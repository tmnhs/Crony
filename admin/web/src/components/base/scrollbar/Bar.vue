<template>
	<div class="scrollbar__track" ref="track" @mousedown="handleTrackMousedown">
		<div class="scrollbar__thumb" ref="thumb" @mousedown="handleThumbMousedown"></div>
	</div>
</template>

<script>
const BAR_MAP = {
	vertical: {
		name: 'vertical',
		offset: 'offsetTop',
		clientSize: 'clientHeight',
		scroll: 'scrollTop',
		scrollSize: 'scrollHeight',
		thumbSize: 'height',
		client: 'clientY',
		position: 'top',
	},
	horizontal: {
		name: 'horizontal',
		offset: 'offsetLeft',
		clientSize: 'clientWidth',
		scroll: 'scrollLeft',
		scrollSize: 'scrollWidth',
		thumbSize: 'width',
		client: 'clientX',
		position: 'left',
	},
}

export default {
	props: {
		// 垂直滚动/水平滚动
		vertical: {
			type: Boolean,
			default: true,
		},
		// 滚动条粗细
		size: {
			type: Number,
			default: 6,
		},
		// 滚动条颜色
		color: {
			type: String,
			default: '#939393',
		},
		// 内容滚动的距离
		scroll: {
			type: Number,
			default: 0,
		},
		// 可视区大小
		clientSize: {
			type: Number,
			default: 0,
		},
		// 内容真实大小
		scrollSize: {
			type: Number,
			default: 0,
		},
	},
	data() {
		return {
			// 滚动轨道
			track: null,
			// 滚动条
			thumb: null,
			// 滚动条长度
			thumbSize: 0,
			// 滚动速率
			speed: 0,
			// 鼠标按下时，光标离滚动条顶部的距离加上滚动区离浏览器顶部的距离
			record: 0,
			bar: BAR_MAP[this.vertical ? 'vertical' : 'horizontal'],
		}
	},
	watch: {
		// 内容区滚动的时候带动滚动条移动
		scroll(value) {
			this.setThumbPosition()
		},
		clientSize() {
			this.handleUpdate()
		},
		scrollSize() {
			this.handleUpdate()
		},
	},
	mounted() {
		this.track = this.$refs.track
		this.thumb = this.$refs.thumb
		this.initTrackStyle()
		this.initThumbStyle()
		this.handleUpdate()
	},
	destroyed() {
		document.removeEventListener('mouseup', this.handleDocumentMouseup)
	},
	methods: {
		// 初始化滚动轨道样式
		initTrackStyle() {
			const { track, size, vertical, bar } = this
			if (vertical) {
				track.style.width = size + 'px'
				track.style.right = '0px'
			} else {
				track.style.height = size + 'px'
				track.style.bottom = '0px'
			}
		},
		// 初始化滚动条样式
		initThumbStyle() {
			const { thumb, color, size, vertical } = this
			thumb.style.backgroundColor = color
			thumb.style.borderRadius = size / 2 + 'px'
			if (vertical) {
				thumb.style.width = size + 'px'
			} else {
				thumb.style.height = size + 'px'
			}
		},
		// 设置滚动条长度
		setThumbSize() {
			const { scrollSize, clientSize, bar } = this
			if (scrollSize > clientSize) {
				this.thumbSize = (clientSize * clientSize) / scrollSize
			} else {
				this.thumbSize = 0
			}
			this.thumb.style[bar.thumbSize] = this.thumbSize + 'px'
		},
		// 设置滚动条位置
		setThumbPosition() {
			const { thumb, scroll, speed, bar } = this
			thumb.style[bar.position] = scroll / speed + 'px'
		},
		// 设置滚动速率
		setScrollSpeed() {
			const { scrollSize, clientSize, thumbSize } = this
			this.speed = (scrollSize - clientSize) / (clientSize - thumbSize)
		},
		handleUpdate() {
			this.setThumbSize()
			this.setScrollSpeed()
			this.setThumbPosition()
		},
		// 控制内容区的滚动
		handleContentScroll(position) {
			let end = position
			const { clientSize, thumbSize, speed, bar } = this
			const max = clientSize - thumbSize
			if (end < 0) {
				// 滚动条到达最顶部
				end = 0
			} else if (end > max) {
				// 滚动条到达最底部
				end = max
			}
			this.$emit('onManualScroll', end * speed, bar.scroll)
		},
		// 拖动滚动条进行滚动
		handleThumbMousedown(event) {
			event.stopPropagation()
			// 原生滚动条不能按住右键进行滚动
			if (event.button === 2) {
				return
			}
			this.record = event[this.bar.client] - this.thumb[this.bar.offset]
			document.addEventListener('mousemove', this.handleDocumentMousemove)
			document.addEventListener('mouseup', this.handleDocumentMouseup)
		},
		handleDocumentMousemove(event) {
			const position = event[this.bar.client] - this.record
			this.handleContentScroll(position)
			// 防止拖动太快选中文字
			window.getSelection().removeAllRanges()
		},
		handleDocumentMouseup() {
			document.removeEventListener('mousemove', this.handleDocumentMousemove)
		},
		// 点击滚动轨道进行滚动
		handleTrackMousedown(event) {
			if (event.button === 2) {
				return
			}
			// 浏览器原生滚动条是点一下，则滚动一定距离，这里直接滚动到点击的位置。
			const position =
				event[this.bar.client] - event.target.getBoundingClientRect()[this.bar.position] - this.thumbSize / 2
			this.handleContentScroll(position)
		},
	},
}
</script>

<style lang="scss" scoped>
.scrollbar__track {
	position: absolute;
	width: 100%;
	height: 100%;

	.scrollbar__thumb {
		position: absolute;
		width: 100%;
		height: 100%;
		cursor: pointer;
		opacity: 0.4;

		&:hover {
			opacity: 0.6;
		}
	}
}
</style>
