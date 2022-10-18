const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const System = () => import(/* webpackChunkName:'system' */ '@/pages/system')
import i18n from  "@/assets/lang";
const route = {
	name: 'System',
	path: '/system',
	component: InnerLayout,
	redirect: '/system/index/',
	meta: {
		title: i18n.t('menu.server'),
		icon: 'data',
		noCache: false,

	},
	children: [
		{
			name: 'system',
			path: '/node/system/:uuid',
			component: System,
			props: true,
			meta: {
				title:  i18n.t('menu.server'),
				hiddenInMenu: true,
				noCache: true,
			},
		},
		{
			name: 'system',
			path: '/system/index/',
			component: System,
			meta: {
				title: i18n.t('menu.server'),
				activePath: '/system',
			},
		}
	],
}

export default route
