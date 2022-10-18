import Vue from 'vue'

const context = require.context('./modules', false, /\.js$/)

context.keys().forEach(path => {
	const directiveName = path
		.split('/')
		.pop()
		.replace(/\.\w+$/, '')
	const directiveConfig = context(path)
	Vue.directive(directiveName, directiveConfig.default)
})
