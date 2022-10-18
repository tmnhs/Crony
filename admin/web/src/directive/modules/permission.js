// DOM级权限控制，若当前的角色不在指令传入的权限数组中，则该DOM元素不渲染。
// TODO  应该采用权限码的方式
import store from '@/store'

export default {
	inserted(el, binding) {
		const roles = binding.value
		if (roles && roles instanceof Array && roles.length > 0) {
			const currentRoles = store.getters.userInfo.roles
			if (!currentRoles.some(role => roles.includes(role))) {
				el.parentNode && el.parentNode.removeChild(el)
			}
		} else {
			throw new Error(`需要传入像这样格式的指令 v-permission="['admin','editor']"`)
		}
	},
}
