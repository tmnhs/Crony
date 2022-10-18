const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const Log = () => import(/* webpackChunkName:'Log' */ '@/pages/log')
import i18n from  "@/assets/lang";
const route = {
	name: 'Log',
	path: '/log',
	component: InnerLayout,
	redirect: '/log/index/',
	meta: {
		title: i18n.t('menu.log'),
		icon: 'log',
		noCache: false,

	},
	children: [
		{
			name: 'JobLog',
			path: '/job/log/:jid',
			component: Log,
			props: true,
			meta: {
				title: i18n.t('menu.log'),
				hiddenInMenu: true,
				noCache: true,
			},
		},
		{
			name: 'Log',
			path: '/log/index/',
			component: Log,
			meta: {
				title: i18n.t('menu.log'),
				activePath: '/log',
			},
		}
	],
}

export default route
