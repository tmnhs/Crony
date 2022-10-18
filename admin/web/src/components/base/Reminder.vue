<template>
	<transition name="reminder-fade">
		<div class="reminder" :class="typeClass" v-show="visible">
			<i class="reminder__icon" :class="iconClass"></i>
			<span class="reminder__content" v-text="message"></span>
			<i class="reminder__closeBtn el-icon-close" @click="close"></i>
		</div>
	</transition>
</template>
<script>
export default {
	data() {
		return {
			visible: false,
			type: 'info',
			message: '',
			duration: 3000,
			onClose: null,
			timer: null,
		}
	},
	computed: {
		typeClass() {
			return `reminder--${this.type}`
		},
		iconClass() {
			const TYPE_CLASSES_MAP = {
				info: 'el-icon-info',
				success: 'el-icon-success',
				warning: 'el-icon-warning',
				error: 'el-icon-error',
				loading: 'el-icon-loading',
			}
			return TYPE_CLASSES_MAP[this.type]
		},
	},
	mounted() {
		this.startTimer()
		document.addEventListener('keydown', this.keydown)
	},
	destroyed() {
		// 销毁实例的时候必须移除对docoment的事件监听，不然会造成重复监听
		document.removeEventListener('keydown', this.keydown)
	},
	methods: {
		close() {
			// 防止重复点击关闭，点一次之后动画要经过一段时间才结束，此时组件还未完全消失还可以再点击。
			if (!this.visible) return
			clearTimeout(this.timer)
			this.visible = false
			//监听过渡结束事件，也可以使用transition组件的atfer-leave钩子
			this.$el.addEventListener('transitionend', () => {
				// 销毁实例是为了触发destroyed钩子
				this.$destroy()
				// 将组件从DOM中移除
				this.$el.parentNode.removeChild(this.$el)
			})
			// 关闭之后的回调
			if (typeof this.onClose === 'function') {
				// 参数为当前reminder实例
				this.onClose(this)
			}
		},
		// 自动关闭
		startTimer() {
			// duration为0，不会自动关闭
			if (this.duration !== 0) {
				this.timer = setTimeout(this.close, this.duration)
			}
		},
		// 监听ESC键
		keydown(event) {
			if (event.keyCode === 27) {
				this.close()
			}
		},
	},
}
</script>
<style lang="scss" scoped>
.reminder {
	min-width: 380px;
	box-sizing: border-box;
	position: fixed;
	left: 50%;
	top: 50%;
	transform: translate(-50%, 50%);
	padding: 16px;
	display: flex;
	align-items: center;
	border-width: 1px;
	border-style: solid;
	border-radius: 4px;
	z-index: 1000;

	.reminder__icon {
		margin-right: 10px;
	}

	.reminder__content {
		padding-right: 16px;
	}

	.reminder__closeBtn {
		position: absolute;
		top: 50%;
		right: 15px;
		transform: translateY(-50%);
		font-size: 16px;
		color: #c0c4cc;
		cursor: pointer;

		&:hover {
			color: #aaa;
		}
	}
}

.reminder--info {
	background-color: #edf2fc;
	border-color: #ebeef5;
	color: #909399;
}

.reminder--success {
	background-color: #f0f9eb;
	border-color: #e1f3d8;
	color: #67c23a;
}

.reminder--warning {
	background-color: #fdf6ec;
	border-color: #faecd8;
	color: #e6a23c;
}

.reminder--error {
	background-color: #fef0f0;
	border-color: #fde2e2;
	color: #f56c6c;
}

.reminder--loading {
	background-color: mix(#fff, #409eff, 90%);
	border-color: mix(#fff, #409eff, 80%);
	color: #409eff;
}

.reminder-fade-enter,
.reminder-fade-leave-to {
	opacity: 0;
}

.reminder-fade-enter-active,
.reminder-fade-leave-active {
	transition: opacity 0.5s;
}
</style>
