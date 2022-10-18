<template>
	<!-- 导航标签，同时负责缓存页面 -->
	<scrollbar>
		<div class="nav-tag">
			<div
				class="nav-tag__item"
				:class="{ active: isTagActive(tag.path) }"
				v-for="(tag, index) in tagList"
				:key="tag.path"
				@dblclick="handleClose(index)"
			>
				<router-link class="link" :to="tag.path">{{ tag.title }}</router-link>
				<i class="el-icon-circle-close close" @click="handleClose(index)"></i>
			</div>
		</div>
	</scrollbar>
</template>

<script>
import Scrollbar from '@/components/base/scrollbar'
import eventBus from '@/utils/event-bus'

export default {
	components: {
		Scrollbar,
	},
	data() {
		return {
			tagList: [],
		}
	},
	watch: {
		$route(route) {
			this.addTag(route)
		},
		tagList() {
			// 缓存页面
			const cachePages = this.tagList.filter(item => !item.noCache).map(item => item.name)
			// console.log(cachePages);
			this.$store.commit('SET_CACHE_PAGES', cachePages)
		},
	},
	created() {
		this.addTag(this.$route)
		eventBus.$on('closeTag', path => {
			const tagIndex = this.tagList.findIndex(tag => tag.path === path)
			this.tagList.splice(tagIndex, 1)
		})
	},
	methods: {
		isTagActive(path) {
			return this.$route.path === path
		},
		// 添加标签
		addTag(route) {
			let tagExist = this.tagList.some(item => item.path == route.path)
			// 路由名称存在才在tag栏中显示，比如Reload路由没设置name，就不显示tag。
			if (!tagExist&&route.name) {
				this.tagList.push({
					title: route.meta.title,
					path: route.path,
					name: route.name,
					noCache: route.meta.noCache,
				})
			}
		},
		// 关闭标签
		handleClose(index) {
			// taglist中只有首页时不关闭首页
			if (this.tagList.length === 1 && this.tagList[0].path === '/dashboard') return false
			const delTag = this.tagList.splice(index, 1)[0]
			//判断关闭的tag的下一个tag存不存在，存在就跳到下一个tag，不存在就跳到上一个tag。
			const toTag = this.tagList[index] ? this.tagList[index] : this.tagList[index - 1]
			if (toTag) {
				// 只有关闭的是当前页面才进行跳转
				delTag.path == this.$route.path && this.$router.push(toTag.path)
			} else {
				// 如果上一个tag也不存在就跳到首页(也就是关闭完所有tag了)
				this.$router.push('/dashboard')
			}
		},
	},
}
</script>

<style lang="scss" scoped>
.scrollbar-wrap {
	flex: none;
	height: 40px;
	background-color: #fff;
	box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 1px 0 rgba(0, 0, 0, 0.1), 0 2px 1px -1px rgba(0, 0, 0, 0.1);
	z-index: 10;
}

.nav-tag {
	display: flex;
	height: 40px;
	box-sizing: border-box;
	padding: 5px 10px;
	white-space: nowrap;

	.nav-tag__item {
		padding: 5px 10px;
		border-radius: 4px;
		border: $base-border;
		background-color: var(--theme-white__9);
		color: #409eff;

		& + .nav-tag__item {
			margin-left: 10px;
		}

		.link {
			display: inline-block;
			font-size: 12px;
			margin-right: 5px;

			&:hover {
				color: $black9 !important;
			}
		}

		.close {
			color: var(--theme);
			cursor: pointer;

			&:hover {
				color: $black9;
			}
		}

		&.active {
			background-color: var(--theme);

			.link,
			.close {
				color: #fff;

				&:hover {
					color: $black9;
				}
			}
		}
	}
}
</style>
