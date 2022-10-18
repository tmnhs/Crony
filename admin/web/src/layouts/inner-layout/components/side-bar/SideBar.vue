<template>
	<div class="side-bar">
		<logo></logo>

		<scrollbar>
			<el-menu
				:default-active="activePath"
				:collapse="sideCollapse"
				:unique-opened="false"
				router
				background-color="#545c64"
				text-color="#fff"
			>
				<menu-item v-for="menu in menuList" :key="menu.path" :config="menu"></menu-item>
			</el-menu>
		</scrollbar>
	</div>
</template>

<script>
import _ from 'lodash'
import Scrollbar from '@/components/base/scrollbar'
import Logo from './Logo'
import MenuItem from './MenuItem'

export default {
	components: {
		Scrollbar,
		Logo,
		MenuItem,
	},
	data() {
		return {
			menuList: [],
		}
	},
	computed: {
		sideCollapse() {
			return this.$store.getters.sideCollapse
		},
		routeMap() {
			return this.$store.getters.routeMap
		},
		activePath() {
			const { meta, path } = this.$route
			return meta.activePath ? meta.activePath : path
		},
	},
	created() {
		const menuRouteMap = this.getMenuRouteMap(this.routeMap)
		this.menuList = menuRouteMap.map(route => this.getMenuItem(route))
	},
	methods: {
		// 获取需要在侧边菜单显示的路由表
		getMenuRouteMap(routes) {
			return _.cloneDeep(routes).filter(route => {
				if (route.meta.hiddenInMenu) {
					return false
				} else {
					if (route.children) {
						const children = this.getMenuRouteMap(route.children)
						route.children = children.length > 0 ? children : undefined
					}
					return true
				}
			})
		},
		getMenuItem(route) {
			// children不存在代表是最后一级路由，只有一个children代表只有第一级路由
			return {
				title: route.meta.title,
				icon: route.meta.icon,
				path: route.path,
				children:
					!route.children || route.children.length === 1
						? undefined
						: route.children.map(childrenRoute => this.getMenuItem(childrenRoute)),
			}
		},
	},
}
</script>

<style lang="scss" scoped>
.side-bar {
	background-color: #545c64;

	.scrollbar-wrap {
		height: calc(100% - 48px);
	}

	.el-menu {
		height: 100%;
		border-right: none;

		&:not(.el-menu--collapse) {
			width: 200px;
		}
	}
}
</style>
