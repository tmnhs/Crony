<template>
	<el-drawer
		class="base-drawer"
		:custom-class="customClass"
		:title="title"
		:visible="visible"
		:direction="direction"
		:size="size"
		:destroy-on-close="destroyOnClose"
		@close="handleClose"
	>
		<div slot="title" v-if="!title">
			<slot name="title">{{ title }}</slot>
		</div>

		<slot></slot>

		<div v-if="showFooter" class="footer" :style="{ width: size }">
			<slot name="footer">
				<el-button v-if="showCancel" @click="handleClose">
					{{ cancelText || '取消' }}
				</el-button>
				<el-button
					v-if="showConfirm"
					:type="type"
					:disabled="confirmDisabled"
					:loading="confirmLoading"
					@click="handleConfirm"
				>
					{{ confirmText || '确认' }}
				</el-button>
			</slot>
		</div>
	</el-drawer>
</template>

<script>
/**
 * 侧边抽屉
 */
export default {
	props: {
		customClass: String,
		title: String,
		direction: String,
		size: {
			type: String,
			default: '400px',
		},
		visible: Boolean,
		type: {
			default: 'primary',
			validator: function (value) {
				return ['primary', 'danger', 'warning', 'info', 'success'].includes(value)
			},
		},
		confirmText: String,
		cancelText: String,
		confirmLoading: Boolean,
		confirmDisabled: Boolean,
		showFooter: {
			type: Boolean,
			default: true,
		},
		showConfirm: {
			type: Boolean,
			default: true,
		},
		showCancel: {
			type: Boolean,
			default: true,
		},
		destroyOnClose: {
			type: Boolean,
			default: true,
		},
	},
	data() {
		return {}
	},
	methods: {
		handleConfirm() {
			this.$emit('confirm')
		},
		handleClose() {
			this.$emit('close')
		},
	},
}
</script>

<style lang="scss" scoped>
.base-drawer {
	/deep/ .el-drawer__header {
		height: 50px;
		padding: 0 0 0 20px;
	}

	/deep/ .el-drawer__close-btn {
		width: 50px;
		height: 50px;
		font-size: 24px;
		color: #909399;

		&:hover {
			color: #373737;
		}

		&:active {
			background-color: #efefef;
		}
	}

	/deep/ .el-drawer__body {
		padding: 20px;
		padding-bottom: 60px;
		overflow: auto;
		position: relative;
	}

	.footer {
		position: fixed;
		bottom: 0;
		height: 60px;
		display: flex;
		align-items: center;
		justify-content: center;
		background-color: #fff;
		transform: translateX(-20px);
	}
}
</style>
