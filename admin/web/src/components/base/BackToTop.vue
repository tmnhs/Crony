<template>
	<transition name="fade">
		<div
			class="rocket"
			ref="rocket"
			v-show="visible"
			@click.stop="clickRocket"
			@mouseenter="enterRocket"
			@mouseleave="leaveRocket"
		></div>
	</transition>
</template>

<script>
/**
 * 滚回顶部
 */
import { scroll } from '@/utils/core'
import Animate from '@/utils/animate'

export default {
	props: {
		//滚动容器
		container: {
			type: String,
			default: '.inner-layout__page',
		},
		// 页面向下滚动到多少距离才显示火箭
		visibleHeight: {
			type: Number,
			default: 100,
		},
		// 火箭的位置
		position: {
			type: Object,
			default() {
				return {
					bottom: '20px',
					right: '30px',
				}
			},
		},
		// 火箭滚动所用的时间
		duration: {
			type: Number,
			default: 500,
		},
	},
	data() {
		return {
			isClickRocket: false,
			visible: false,
			scrollElem: null,
		}
	},
	mounted() {
		this.scrollElem = document.querySelector(this.container)
		this.scrollElem.addEventListener('scroll', this.handleRocketShow)
		const rocket = this.$refs.rocket
		rocket.style.bottom = this.position.bottom
		rocket.style.right = this.position.right
	},
	destroyed() {
		this.scrollElem.removeEventListener('scroll', this.handleRocketShow)
	},
	methods: {
		// 当页面向下滚动到移动距离时火箭才显示
		handleRocketShow() {
			if (this.scrollElem.scrollTop > this.visibleHeight) {
				this.visible = true
			} else {
				// 如果没有点击火箭，火箭才不显示，不然当火箭向上移动到visibleHeight的位置时就会消失
				if (this.isClickRocket === false) {
					this.visible = false
				}
			}
		},
		// 火箭和页面的运动都结束之后执行的回调
		callback() {
			this.isClickRocket = false
			this.visible = false
			this.changeRocketImagePosition('-31px -15px')
			// 火箭不能在向上移动到目标位置之后就立刻定位到原来的位置，因为它消失的过程有个渐进的动画，完全隐藏之后再进行定位。
			setTimeout(() => {
				this.$refs.rocket.style.bottom = this.position.bottom
			}, 500)
		},
		// 点击火箭的时候，一方面页面先快后慢滚动到顶部，一方面火箭先慢后快移动直到消失。
		clickRocket(event) {
			this.isClickRocket = true
			this.changeRocketImagePosition('-204px -15px')
			const callback = () => {
				this.isClickRocket = false
				this.visible = false
				this.changeRocketImagePosition('-31px -15px')
				// 火箭不能在向上移动到目标位置之后就立刻定位到原来的位置，因为它消失的过程有个渐进的动画，完全隐藏之后再进行定位。
				setTimeout(() => {
					this.$refs.rocket.style.bottom = this.position.bottom
				}, 500)
			}
			scroll(this.scrollElem, 0, 800)
			const animate = new Animate(this.$refs.rocket)
			//设置火箭移动花费的时间比页面要慢点，然后在火箭运动结束后执行回调
			animate.start('bottom', window.innerHeight, 850, 'easeIn', callback)
		},
		enterRocket() {
			this.changeRocketImagePosition('-117px -15px')
		},
		leaveRocket() {
			if (this.isClickRocket == false) {
				this.changeRocketImagePosition('-31px -15px')
			}
		},
		// 切换火箭的图片
		changeRocketImagePosition(position) {
			this.$refs.rocket.style.backgroundPosition = position
		},
	},
}
</script>

<style scoped>
.fade-enter,
.fade-leave-to {
	opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
	transition: opacity 0.5s;
}

.rocket {
	position: fixed;
	width: 31px;
	height: 76px;
	background: url(~@/assets/images/common/rocket.png) no-repeat -31px -15px;
	cursor: pointer;
	z-index: 10000;
}
</style>
