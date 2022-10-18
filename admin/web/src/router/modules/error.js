const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const Forbidden = () => import(/* webpackChunkName:'forbidden' */ '@/pages/error/forbbiden')
const NotFound = () => import(/* webpackChunkName:'notFound' */ '@/pages/error/not-found')
import i18n from  "@/assets/lang";
const route = {
	name: 'Error',
	path: '/error',
	component: InnerLayout,
	meta: {
		title: i18n.t('menu.error'),
		icon: 'error',
		hiddenInMenu:true,
	},
	children: [
		{
			name: 'Forbidden',
			path: '/error/forbidden',
			component: Forbidden,
			meta: {
				title: '403',
				noCache: true,
				hiddenInMenu:true,
			},
		},
		{
			name: 'NotFound',
			path: '/error/notFound',
			component: NotFound,
			meta: {
				title: '404',
				noCache: true,
				hiddenInMenu:true,

			},
		},
	],
}

export default route
