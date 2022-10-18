const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const Dashboard = () => import(/* webpackChunkName:'dashboard' */ '@/pages/dashboard')
import i18n from  "@/assets/lang";
const route = {
	name: 'Dashboard',
	path: '/dashboard',
	component: InnerLayout,
	redirect: '/dashboard/index',
	meta: {
		title: i18n.t('menu.home'),
		icon: 'home',
	},
	children: [
		{
			path: '/dashboard/index',
			component: Dashboard,
			meta: {
				title: i18n.t('menu.home'),
				activePath: '/dashboard',
			},
		},
	],
}

export default route
