export default {
	token: state => state.account.token,
	userInfo: state => state.account.userInfo,
	routeMap: state => state.account.routeMap,
	cachePages: state => state.cache.cachePages,
	theme: state => state.setting.theme,
	sideCollapse: state => state.setting.sideCollapse,
	tagVisible: state => state.setting.tagVisible,
	style: state => state.setting.style,
	size: state => state.setting.size,
	language:state => state.language.language,
}
