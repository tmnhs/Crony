<template>
	<el-dialog
		class="base-dialog"
		:custom-class="customClass"
		:title="title"
		:visible="visible"
		:width="width"
		:top="top"
		:destroy-on-close="destroyOnClose"
		@close="handleClose"
	>
		<div slot="title" v-if="!title">
			<slot name="title"></slot>
		</div>

		<slot></slot>

		<div v-if="showFooter" class="footer" slot="footer">
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
	</el-dialog>
</template>

<script>
/**
 * 模态窗
 */
export default {
	props: {
		customClass: String,
		title: String,
		visible: Boolean,
		type: {
			default: 'primary',
			validator: function (value) {
				return ['primary', 'danger', 'warning', 'info', 'success'].includes(value)
			},
		},
		top: String,
		width: {
			type: String,
			default: '500px',
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
.base-dialog {
	/deep/.el-dialog__header {
		display: flex;
		align-items: center;
		height: 50px;
		padding: 0 0 0 20px;
	}

	/deep/ .el-dialog__headerbtn {
		width: 50px;
		height: 50px;
		position: static;
		margin-left: auto;
		font-size: 24px;

		&:active {
			background-color: #efefef;
		}
	}

	/deep/ .el-dialog__body {
		padding: 20px;
	}

	/deep/ .el-dialog__footer {
		text-align: center;
	}
}
</style>
