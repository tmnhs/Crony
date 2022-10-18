<template>
	<svg class="base-icon" :class="iconClass" :style="iconStyle" aria-hidden="true" @click="handleClick">
		<use :xlink:href="iconName" :stroke-width="strokeWidth"></use>
	</svg>
</template>

<script>
/**
 * svg图标
 * 图标库使用的是iconpark: https://iconpark.oceanengine.com。
 * 如果要使用单色图标，将svg图标代码中的stroke属性修改为currentColor。
 * 如果需要自定义图标线条的粗细，将svg代码中的stroke-width属性删除或修改为想要的粗细。
 * vscode安装SVG插件，方便修改svg文件的时候对svg代码格式化。
 * vscode安装Svg Preview插件预览svg。
 * win平台安装SVG Explorer Extension，查看模式设为“中图标”，可直接在系统文件夹中查看所有图标。
 */

export default {
	props: {
		// 图标名称，就是svg文件的名称
		name: {
			type: String,
			default: '',
		},
		// 样式
		iconClass: {
			type: String,
			default: '',
		},
		// 图标大小
		size: {
			type: [String, Number],
			default: 18,
		},
		//  图标颜色
		color: {
			type: String,
		},
		// 线条的粗细
		strokeWidth: {
			type: [String, Number],
			default: 4,
		},
	},
	computed: {
		iconName() {
			return '#icon-' + this.name
		},
		iconStyle() {
			return {
				fontSize: typeof this.size === 'string' ? this.size : `${this.size}px`,
				color: this.color,
			}
		},
	},
	methods: {
		handleClick() {
			this.$emit('click')
		},
	},
}
</script>

<style scoped>
.base-icon {
	width: 1em;
	height: 1em;
	vertical-align: -0.15em;
	overflow: hidden;
	fill: currentColor;
	stroke: currentColor;
}
</style>
