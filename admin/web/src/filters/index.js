import Vue from 'vue'

const context = require.context('./modules', false, /\.js$/)

context.keys().forEach(path => {
	const filterName = path
		.split('/')
		.pop()
		.replace(/\.\w+$/, '')
	const filterConfig = context(path)
	Vue.filter(filterName, filterConfig.default)
})
