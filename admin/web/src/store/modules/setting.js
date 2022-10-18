import uiConfig from '@/config/ui'
import { localMng } from '@/utils/storage-mng'

export default  {
	state: {
		theme: localMng.getItem('theme') || uiConfig.theme,
		sideCollapse: localMng.getItem('sideCollapse') || uiConfig.sideCollapse,
		tagVisible: localMng.getItem('tagVisible') || uiConfig.tagVisible,
		style: localMng.getItem('style') || uiConfig.style,
		size: localMng.getItem('size') || uiConfig.size,
	},
	mutations: {
		SET_THEME(state, theme) {
			state.theme = theme
			localMng.setItem('theme', theme)
		},
		SET_SIDE_COLLAPSE(state, collapse) {
			state.sideCollapse = collapse
			localMng.setItem('sideCollapse', collapse)
		},
		SET_TAG_VISIBLE(state, visible) {
			state.tagVisible = visible
			localMng.setItem('tagVisible', visible)
		},
		SET_STYLE(state, style) {
			state.style = style
			localMng.setItem('style', style)
		},
		SET_SIZE(state, size) {
			state.size = size
			localMng.setItem('size', size)
		},
	},
}

