const InnerLayout = () => import(/* webpackChunkName:'innerLayout' */ '@/layouts/inner-layout')
const JobList = () => import(/* webpackChunkName:'jobList' */ '@/pages/job/list')
const JobEdit = () => import(/* webpackChunkName:'jobEdit' */ '@/pages/job/edit')
import i18n from  "@/assets/lang";

const route = {
	name: 'Job',
	path: '/job',
	component: InnerLayout,
	redirect: '/job/list',
	meta: {
		title: i18n.t('menu.job'),
		icon: 'chart',
	},
	children: [
		{
			name: 'JobList',
			path: '/job/list',
			component: JobList,
			meta: {
				title: i18n.t('menu.job'),
				activePath: '/job',
				noCache: false,
			},
		},
		{
			name: 'NodeJobList',
			path: '/node/job/:uuid',
			component: JobList,
			props: true,
			meta: {
				title: i18n.t('menu.node_job'),
				hiddenInMenu: true,
				noCache: true,
			},
		},
		{
			name: 'JobAdd',
			path: '/job/add',
			component: JobEdit,
			meta: {
				title: i18n.t('menu.add_job'),
				hiddenInMenu: true,
				noCache: true,
			},
		},
		{
			name: 'JobEdit',
			path: '/job/edit/:job_id',
			component: JobEdit,
			meta: {
				title: i18n.t('menu.edit_job'),
				hiddenInMenu: true,
				noCache: true,
			},
			props: true,
			beforeEnter: (to, from, next) => {
				to.meta.title = i18n.t('menu.edit_job') + '-' + to.params.job_id
				next()
			},
		},
	],
}

export default route
