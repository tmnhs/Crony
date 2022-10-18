<template>
	<div class="inner-layout">
		<side-bar></side-bar>

		<div class="inner-layout__main">
			<header-bar></header-bar>
			<nav-tag v-show="tagVisible"></nav-tag>

			<div class="inner-layout__page">
				<transition name="page" mode="out-in">
					<keep-alive :include="cachePages">
						<router-view :key="key"></router-view>
					</keep-alive>
				</transition>
			</div>
		</div>
	</div>
</template>

<script>
/**
 * 内部布局(已登录)
 */
import { HeaderBar, SideBar, NavTag } from './components'
import { mapGetters } from 'vuex'

export default {
	components: { HeaderBar, SideBar, NavTag },
	computed: {
		...mapGetters(['tagVisible', 'cachePages']),
		// 创建文章以及所有的编辑文章都是共用的ArticleEdit组件,所以对应的路由不能使用keep-alive缓存。
		// 同时它们使用的也是动态路由，在这些路由之间切换时也需要响应路由参数的变化，当路由参数变化时需要重新渲染。
		key() {
			return this.$route.path
		},
	},
}
</script>

<style lang="scss" scoped>
.inner-layout {
	display: flex;
	height: 100vh;
	overflow-y: hidden;

	.inner-layout__main {
		flex: auto;
		display: flex;
		flex-direction: column;
		overflow-x: hidden;
		background-color: #f0f2f5;

		.inner-layout__page {
			flex: auto;
			position: relative;
			box-sizing: border-box;
			overflow-x: hidden;
			overflow-y: auto;
			padding: 1em;

			.page {
				&-enter {
					opacity: 0;
					transform: translateX(-30px);
				}

				&-leave-to {
					opacity: 0;
					transform: translateX(30px);
				}

				&-enter-active,
				&-leave-active {
					transition: all 0.3s ease;
				}
			}
		}
	}
}
</style>
