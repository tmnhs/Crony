<template>
	<div>
		<el-upload
			class="drag-upload"
			:action="action"
			:name="name"
			drag
			multiple
			:file-list="value"
			:before-upload="beforeUpload"
			:before-remove="beforeRemove"
			:on-success="handleSuccess"
			:on-remove="handleRemove"
			:on-preview="handlePreview"
		>
			<i class="drag-upload__icon" :class="loading ? 'el-icon-loading' : 'el-icon-upload '"></i>
			<p class="drag-upload__text">点击或直接将文件拖到此处上传</p>
			<p class="drag-upload__tip">文件大小不能超过{{ sizeLimit }}MB！{{ tip }}</p>
		</el-upload>
	</div>
</template>

<script>
/**
 * 文件拖拽上传
 */
export default {
	props: {
		// 文件列表
		value: {
			type: Array,
			default() {
				return []
			},
		},
		//上传地址
		action: {
			required: true,
			type: String,
			default: '',
		},
		// 对应inpu控件的name属性，后端依据这个字段获取文件。
		name: {
			type: String,
			default: 'file',
		},

		// 文件的大小限制,单位为MB
		sizeLimit: {
			type: Number,
			default: 10,
		},
		// 提示信息
		tip: {
			type: String,
			default: '',
		},
	},
	data() {
		return {
			loading: false,
		}
	},
	methods: {
		beforeUpload(file) {
			const limit = file.size / 1024 / 1024 < this.sizeLimit
			if (!limit) {
				this.$message.error(`上传的文件小不能超过 ${this.sizeLimit}MB!`)
			}
			if (limit) {
				this.loading = true
			}
			return limit
		},
		beforeRemove(file, fileList) {
			return this.$confirm(`确定删除“${file.name}”？`)
		},
		handleSuccess(res, file, fileList) {
			this.loading = false
			//根据实际开发情况处理响应
			if (true) {
				this.$emit('input', fileList)
			} else {
				this.$message.error(res.message || '上传失败')
			}
		},
		handleRemove(file, fileList) {
			this.$emit('input', fileList)
		},
		handlePreview(file) {
			window.open(file.url)
		},
	},
}
</script>

<style lang="scss" scoped>
.drag-upload {
	.drag-upload__icon {
		font-size: 40px;
		line-height: 40px;
		color: var(--theme);
		margin: 0;
	}

	.drag-upload__text {
		line-height: 20px;
		margin-bottom: 6px;
	}

	.drag-upload__tip {
		font-size: 12px;
		line-height: 20px;
		color: $black9;
	}

	/deep/.el-upload {
		width: 100%;
	}

	/deep/.el-upload-dragger {
		width: 100%;
		min-height: 140px;
		height: 100%;
		padding: 20px 1em;
	}
}
</style>
