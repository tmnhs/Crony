<template>
	<div @mouseenter="checkOverflow">
		<el-tooltip v-if="isOverFlow && showTooltip" class="truncate">
			<template slot="content">
				<template v-if="title">{{ title }}</template>
				<slot v-else></slot>
			</template>

			<div>
				<slot></slot>
			</div>
		</el-tooltip>

		<div v-else class="truncate">
			<slot></slot>
		</div>
	</div>
</template>

<script>
export default {
	props: {
		title: String,
		showTooltip: {
			type: Boolean,
			default: true,
		},
	},
	data() {
		return {
			isOverFlow: false,
		}
	},
	methods: {
		checkOverflow(event) {
			const wrap = event.target
			const range = document.createRange()
			range.setStart(wrap, 0)
			range.setEnd(wrap, wrap.childNodes.length)
			const rangeWidth = range.getBoundingClientRect().width
			const isOverFlow = rangeWidth > wrap.clientWidth
			console.log(rangeWidth, wrap.clientWidth)
			this.isOverFlow = isOverFlow
		},
	},
}
</script>

<style scoped>
.truncate {
	overflow: hidden;
	text-overflow: ellipsis;
	white-space: nowrap;
}
</style>
