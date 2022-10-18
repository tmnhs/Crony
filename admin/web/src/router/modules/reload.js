const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const Reload = () => import(/* webpackChunkName:'reload' */ '@/pages/reload')

const route = {
	name: 'Reload',
	path: '/reload',
	component: InnerLayout,
	redirect: '/reload/index',
	meta: {
		hiddenInMenu: true,
	},
	children: [
		{
			name: 'Reload',
			path: '/reload/index',
			component: Reload,
			meta: {
				title: '',
			},
		},
	],
}

export default route
