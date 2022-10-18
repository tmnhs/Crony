const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const NotFound = () => import(/* webpackChunkName:'notFound' */ '@/pages/error/not-found')

// 必须加在所有路由最后。前边的路由都未匹配的时候，才匹配到404
const route = {
	name:"404",
	path: '*',
	component: InnerLayout,
	redirect: '/error/404',
	meta: {
		hiddenInMenu: true,
	},
	children: [
		{
			name: '404',
			path: '/error/404',
			component: NotFound,
		},
	],
}

export default route
