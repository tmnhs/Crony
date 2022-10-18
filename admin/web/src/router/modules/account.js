const OuterLayout = () => import(/* webpackChunkName:'outerLayout' */ '@/layouts/outer-layout')
const Login = () => import(/* webpackChunkName:'login' */ '@/pages/account/login')

const route = {
	name: 'Account',
	path: '/account',
	component: OuterLayout,
	meta: {
		hiddenInMenu: true,
	},
	children: [
		{
			name: 'Login',
			path: '/account/login',
			component: Login,
			meta: {
				title: 'login',
			},
		},
	
	],
}

export default route
