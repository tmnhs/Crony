import Vue from 'vue'
import Vuex from 'vuex'
import _ from 'lodash'
import getters from './getters'

const modules = {}
const context = require.context('./modules', false, /\.js$/)
context.keys().forEach(path => {
	const moduleName = _.camelCase(
		path
			.split('/')
			.pop()
			.replace(/\.\w+$/, '')
	)
	const moduleConfig = context(path)
	modules[moduleName] = moduleConfig.default || moduleConfig
})

Vue.use(Vuex)

const store = new Vuex.Store({
	getters,
	modules,
})

export default store
