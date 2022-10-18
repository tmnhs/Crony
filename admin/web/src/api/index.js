import _ from 'lodash'

const api = {}
const apiContext = require.context('./modules', false, /\.js$/)
apiContext.keys().forEach(path => {
	const moduleName = _.camelCase(
		path
			.split('/')
			.pop()
			.replace(/\.\w+$/, '')
	)
	api[moduleName] = apiContext(path)
})

export default api
