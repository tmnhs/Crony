const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const scriptList = () => import(/* webpackChunkName:'scriptList' */ '@/pages/script/list')
const scriptEdit = () => import(/* webpackChunkName:'scriptEdit' */ '@/pages/script/edit')
const System = () => import(/* webpackChunkName:'system' */ '@/pages/system')
import i18n from  "@/assets/lang";

const route = {
	name: 'Setting',
	path: '/script',
	component: InnerLayout,
	redirect: '/script/list',
	meta: {
		title: '设置',
		icon: 'setting',
	},
	children: [
		{
			name: 'scriptList',
			path: '/script/list',
			component: scriptList,
			meta: {
				title: i18n.t('menu.script'),
				activePath: '/script',
				noCache: false,
				icon: 'script',
			},
		},

		{
			name: 'scriptAdd',
			path: '/script/add',
			component: scriptEdit,
			meta: {
				title: i18n.t('add.script'),
				hiddenInMenu: true,
				noCache: true,
			},
		},
		{
			name: 'scriptEdit',
			path: '/script/edit/:script_id',
			component: scriptEdit,
			meta: {
				title: i18n.t('menu.edit_script'),
				hiddenInMenu: true,
				noCache: true,
			},
			props: true,
			beforeEnter: (to, from, next) => {
				to.meta.title = i18n.t('menu.edit_script') + '-' + to.params.script_id
				next()
			},
		},
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
				icon: 'node',
			},
		}
	],
}

export default route
