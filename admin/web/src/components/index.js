// 注册全局组件
import Vue from 'vue'
import _ from 'lodash'

const context = require.context('./base', true, /\.vue$/)

context.keys().forEach(path => {
	// 获取组件的PascalCase命名
	const componentName = _.upperFirst(
		_.camelCase(
			path
				.split('/')
				.pop()
				.replace(/\.\w+$/, '')
		)
	)
	const componentConfig = context(path)
	Vue.component(componentName, componentConfig.default)
})
