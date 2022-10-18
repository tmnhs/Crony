<template>
	<div class="scrollbar-wrap">
		<div class="scrollbar__content" ref="content" @scroll="handleAutoScroll" @wheel="handleWheel">
			<div ref="resize">
				<slot />
			</div>
		</div>

		<bar
			ref="vertical"
			:size="size"
			:color="color"
			:scroll="scrollTop"
			:clientSize="clientHeight"
			:scrollSize="scrollHeight"
			@onManualScroll="handleManualScroll"
		></bar>

		<bar
			ref="horizontal"
			:vertical="false"
			:size="size"
			:color="color"
			:scroll="scrollLeft"
			:clientSize="clientWidth"
			:scrollSize="scrollWidth"
			@onManualScroll="handleManualScroll"
		></bar>
	</div>
</template>

<script>
/**
 * 滚动条
 */
import Bar from './Bar'
import getScrollbarWidth from './util'
// element直接提供
import { addResizeListener, removeResizeListener } from 'element-ui/src/utils/resize-event'

export default {
	props: ['size', 'color'],
	components: { Bar },
	data() {
		return {
			content: null,
			scrollTop: 0,
			scrollLeft: 0,
			clientHeight: 0,
			clientWidth: 0,
			scrollHeight: 0,
			scrollWidth: 0,
			verticalVisible: true,
			horizontalVisible: true,
		}
	},
	mounted() {
		this.content = this.$refs.content
		this.handleUpdate()
		this.setContentMargin()
		// 没有节流是因为，设置的时间间隔过长会出现卡顿，时间过短跟没节流没什么区别
		// 监听offsetSize的变化
		addResizeListener(this.$refs.content, this.handleUpdate)
		// 监听scrollSize的变化
		addResizeListener(this.$refs.resize, this.handleUpdate)
	},
	destroyed() {
		removeResizeListener(this.$refs.content, this.handleUpdate)
		removeResizeListener(this.$refs.resize, this.handleUpdate)
	},
	methods: {
		getSize() {
			const { scrollHeight, scrollWidth, clientHeight, clientWidth } = this.content
			this.clientWidth = clientWidth
			this.clientHeight = clientHeight
			this.scrollWidth = scrollWidth
			this.scrollHeight = scrollHeight
		},
		getScroll() {
			const { scrollTop, scrollLeft } = this.content
			this.scrollTop = scrollTop
			this.scrollLeft = scrollLeft
		},
		// 获取滚动条显示/隐藏状态
		getVisible() {
			const { scrollHeight, scrollWidth, clientHeight, clientWidth } = this
			this.verticalVisible = scrollHeight > clientHeight
			this.horizontalVisible = scrollWidth > clientWidth
		},
		// 设置负的margin，隐藏原生滚动条
		setContentMargin() {
			const scrollbarWidth = getScrollbarWidth()
			this.content.style.marginRight = this.content.style.marginBottom = -scrollbarWidth + 'px'
		},
		handleUpdate() {
			this.getSize()
			this.getScroll()
			this.getVisible()
		},
		handleManualScroll(value, scroll) {
			this.content[scroll] = value
			this[scroll] = value
		},
		// 垂直滚动监听scroll事件
		handleAutoScroll() {
			if (this.verticalVisible) {
				this.scrollTop = this.content.scrollTop
			}
		},
		// 只存在水平滚动条的时候，水平滚动监听wheel事件
		handleWheel(event) {
			if (!this.verticalVisible && this.horizontalVisible) {
				// 阻止默认事件，不然当页面存在滚动条的时候，会滚动页面
				event.preventDefault()
				// webkit内核的event.deltaY是100，火狐是3，这里指定每次滚动40px；
				let step
				if (event.deltaY > 0) {
					step = 40
				} else {
					step = -40
				}
				this.content.scrollLeft = this.content.scrollLeft + step
				this.scrollLeft = this.content.scrollLeft
			}
		},
	},
}
</script>

<style lang="scss" scoped>
.scrollbar-wrap {
	height: 100%;
	position: relative;
	overflow: hidden;

	.scrollbar__content {
		position: absolute;
		top: 0px;
		left: 0px;
		right: 0px;
		bottom: 0px;
		// 一定要设置成scoll，目的是不管存不存在滚动，都让原生滚动区显示，这样才方便同时设置负的margin-bottom和margin-right
		overflow: scroll;
	}
}
</style>
