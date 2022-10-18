import api from '@/api'
import router, { resetRouter } from '@/router'
import request from '@/utils/request'
import { sessionMng } from '@/utils/storage-mng'

export default {
	state: {
		token: sessionMng.getItem('token'),
		userInfo: {
			id: '',
			name: '',
			avatar: '',
			email: '',
			roles: [],
			routeNames: [],
		},
		routeMap: [],
	},
	mutations: {
		SET_TOKEN(state, token) {
			state.token = token
		},
		SET_USER_INFO(state, userInfo) {
			state.userInfo = userInfo
		},
		SET_ROUTE_MAP(state, routeMap) {
			state.routeMap = routeMap
		},
	},
	actions: {
		// 登录获取token
		async login({ commit }, loginInfo) {
			const username = loginInfo.username.trim()
			const password = loginInfo.password
			const res = await api.account.login({
				username,
				password,
			})
			const token = res.token
			commit('SET_TOKEN', token)
			sessionMng.setItem('token', token)
			request.setHeader({
				Authorization: token,
			})
			router.push('/dashboard')
		},
		// 实际开发token放在请求头的Authorization中
		async getUserInfo({ commit, state }) {
			const res = await api.account.getUserInfo({})
			commit('SET_USER_INFO', res)
			return res
		},
		// 退出
		async logout({ commit }) {

			commit('SET_TOKEN', '')
			commit('SET_USER_INFO', {})
			commit('SET_ROUTE_MAP', [])
			sessionMng.clear()
			request.setHeader({
				Authorization: '',
			})
			resetRouter()
			router.replace('/login')
    },
	},
}
