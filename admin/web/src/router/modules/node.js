const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const Node = () => import(/* webpackChunkName:'Node' */ '@/pages/node')
import i18n from  "@/assets/lang";
const route = {
	name: 'Node',
	path: '/node',
	component: InnerLayout,
	redirect: '/node/index',
	meta: {
		title: i18n.t('menu.node'),
		icon: 'node',
	},
	children: [
		{
			name: 'Node',
			path: '/node/index',
			component: Node,
			meta: {
				title:  i18n.t('menu.node'),
				activePath: '/node',
			},
		},
	],
}

export default route
