const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const User = () => import(/* webpackChunkName:'user' */ '@/pages/user')
import i18n from  "@/assets/lang";
const route = {
	name: 'User',
	path: '/user',
	component: InnerLayout,
	redirect: '/user/index',
	meta: {
		title: i18n.t('menu.user'),
		icon: 'user',
	},
	children: [
		{
			name: 'User',
			path: '/user/index',
			component: User,
			meta: {
				title: i18n.t('menu.user'),
				activePath: '/user',
			},
		},
	],
}

export default route
