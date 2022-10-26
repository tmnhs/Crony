import Vue from 'vue'
import Router from 'vue-router'
import store from '@/store'
import _ from 'lodash'
import { getTreeNodeValue } from '@/utils/core'
import { sessionMng } from '@/utils/storage-mng'
import accountRoute from './modules/account'
import dashboardRoute from './modules/dashboard'
import errorRoute from './modules/error'
import mineRoute from './modules/mine'
import reloadRoute from './modules/reload'
import userRoute from './modules/user'
import jobRoute from './modules/job'
import settingRoute from './modules/setting'
import nodeRoute from './modules/node'
import logRoute from './modules/log'

Vue.use(Router)

const staticRouteMap = [
	{
		path: '/',
		redirect: '/dashboard',
		meta: {
			hiddenInMenu: true,
		},
	},
	accountRoute,
]

const dynamicRouteMap = [
	dashboardRoute,
	mineRoute,
	reloadRoute,
	errorRoute,
	userRoute,
	nodeRoute,
	jobRoute,
	logRoute,
	settingRoute,

	// 必须写在最后。前边的路由都未匹配的时候，才匹配到404
	{
		name: '404',
		path: '*',
		redirect: '/error/notFound',
		meta: {
			hiddenInMenu: true,
		},
	},
]
const routeAdminNames = [
	'Dashboard',
	'User',
	"Log",
	'Mine',
	'Node',
	'Setting',
	'Job',
	'404',
	'Error',
]
const routeNames = [
	'Dashboard',
	'Mine',
	"Log",
	'Node',
	'Job',
	'404',
	'Error',
]
const createRouter = () =>
	new Router({
		// mode: 'history',
		routes: staticRouteMap,
		scrollBehavior(to, from, savedPosition) {
			const innerPage = document.querySelector('.inner-layout__page')
			if (innerPage) {
				innerPage.scrollTo(0, 0)
			}
			return { x: 0, y: 0 }
		},
	})

const router = createRouter()

// 根据路由名称获取可访问的路由表
const filterRouteMap = (routeNames, routeMap) => {
	const acceptedRouteMap = []
	const routes = _.cloneDeep(routeMap)
	routes.forEach(route => {
		// 如果一级路由的名称存在路由权限表中，则它之下的所有子路由都可访问
		if (routeNames.includes(route.name)) {
			acceptedRouteMap.push(route)
		} else {
			// 如果一级路由的名称不在路由权限表中，再看它的哪些子路由名称在路由权限表中
			if (route.children) {
				route.children = filterRouteMap(routeNames, route.children)
				// 如果有子路由可访问，再添加。
				if (route.children.length > 0) {
					acceptedRouteMap.push(route)
				}
			}
		}
	})
	return acceptedRouteMap
}


// 导航守卫
router.beforeEach(async (to, from, next) => {
	const token = sessionMng.getItem('token')
	const outerPaths = getTreeNodeValue([accountRoute], 'path')
	// token不存在(说明没登录),但是路由将要进入系统内部，自动跳到登录页面。
	if (!token && !outerPaths.includes(to.path)) {
		next('/account/login')
	} else {
		// 如果token存在(说明已登录)，但是用户信息不存在，这时应该获取用户信息
		if (token && !store.getters.userInfo.id) {
			const user = await store.dispatch('getUserInfo')
			if(user.role==2){
				//admin
				 const acceptedRouteMap = filterRouteMap(routeAdminNames, dynamicRouteMap)
				// 动态注册路由
				router.addRoutes(acceptedRouteMap)
				store.commit('SET_ROUTE_MAP', [...staticRouteMap, ...acceptedRouteMap])
				//  中断当前导航，重新导航到当前路由。刷新页面之后，会重新注册路由，这样可以确保路由注册完毕后，再进入。
				// replace: true 是为了防止在history中留下之前中断的导航的记录。
			}else {
				const acceptedRouteMap = filterRouteMap(routeNames, dynamicRouteMap)
				// 动态注册路由
				router.addRoutes(acceptedRouteMap)
				store.commit('SET_ROUTE_MAP', [...staticRouteMap, ...acceptedRouteMap])
				//  中断当前导航，重新导航到当前路由。刷新页面之后，会重新注册路由，这样可以确保路由注册完毕后，再进入。
				// replace: true 是为了防止在history中留下之前中断的导航的记录。
			}

			next({ ...to, replace: true })
		} else {
			next()
		}
	}
})

const originalPush = Router.prototype.push
Router.prototype.push = function push(location) {
	return originalPush.call(this, location).catch(err => err)
}

// 退出登录的时候执行，防止重复注册路由
const resetRouter = () => {
	const newRouter = createRouter()
	router.matcher = newRouter.matcher
}

export { resetRouter }

export default router
