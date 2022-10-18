const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const Mine = () => import(/* webpackChunkName:'mine' */ '@/pages/mine')
const Password = () => import(/* webpackChunkName:'mine' */ '@/pages/mine/password')
import i18n from  "@/assets/lang";
const route = {
	name: 'Mine',
	path: '/mine',
	component: InnerLayout,
	redirect: '/mine/index',
	meta: {
		title: i18n.t('user.mine'),
		hiddenInMenu: true,
	},
	children: [
		{
			name: 'Mine',
			path: '/mine/index',
			component: Mine,
			meta: {
				title:  i18n.t('user.mine'),
			},
		},
		{
			name: 'Pw',
			path: '/mine/pw',
			component: Password,
			meta: {
				title:  i18n.t('user.change_pw'),
			},
		},
	],
}

export default route
