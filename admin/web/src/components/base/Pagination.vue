<template>
	<div class="pagination">
		<el-pagination
			:class="`pagination--${position}`"
			:total="total"
			:current-page.sync="currentPage"
			:page-size.sync="pageSizeNum"
			:page-sizes="pageSizes"
			:layout="layout"
			:background="background"
			@size-change="handlePageSizeChange"
			@current-change="handlePageNumberChange"
		>
		</el-pagination>
	</div>
</template>

<script>
/**
 * 分页
 */
export default {
	props: {
		total: {
			required: true,
			type: Number,
		},
		pageNumber: {
			type: Number,
			default: 1,
		},
		pageSize: {
			type: Number,
			default: 10,
		},
		pageSizes: {
			type: Array,
			default() {
				return [10, 20, 30, 50, 1000]
			},
		},
		layout: {
			type: String,
			default: 'total, sizes, prev, pager, next, jumper',
		},
		background: {
			type: Boolean,
			default: true,
		},
		position: {
			validator(value) {
				return ['left', 'center', 'right'].includes(value)
			},
			default: 'center',
		},
	},
	computed: {
		currentPage: {
			get() {
				return this.pageNumber
			},
			set(val) {
				this.$emit('update:pageNumber', val)
			},
		},
		pageSizeNum: {
			get() {
				return this.pageSize
			},
			set(val) {
				this.$emit('update:pageSize', val)
			},
		},
	},
	methods: {
		handlePageSizeChange(val) {
			this.$emit('pagination', val)
		},
		handlePageNumberChange(val) {
			this.$emit('pagination', val)
		},
	},
}
</script>

<style lang="scss">
.pagination {
	padding: 10px;
	height: 50px;

	&--left {
		float: left;
	}

	&--center {
		text-align: center;
	}

	&--right {
		float: right;
	}

	.el-pagination .el-pager {
		.number {
			font-weight: normal;
			color: $black9;
		}

		.el-icon-more {
			border: none !important;
		}
	}

	.el-pagination.is-background .btn-prev,
	.el-pagination.is-background .btn-next,
	.el-pagination.is-background .el-pager li {
		background-color: #fff;
		border: $base-border;
		border-radius: 4px;
	}

	.el-pagination.is-background .el-pager li:not(.disabled).active {
		border-color: var(--theme);
		background-color: #fff;
		color: $black9;
	}
}
</style>
