export default {
	state: {
		// 需要缓存的页面
		cachePages: [],
	},
	mutations: {
		SET_CACHE_PAGES(state, pages) {
			//开发环境不启用,防止热更新失效
			if (process.env.NODE_ENV === 'development') {
				state.cachePages = []
			} else {
				state.cachePages = pages
			}
		},
	},
}
