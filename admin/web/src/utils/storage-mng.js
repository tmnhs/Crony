/**
 * 本地存储的读取往往分散在各个不同的地方，会显得很乱。
 * 使用本地存储的时候推荐统一采用这个封装好的工具，同时在这里记录每个key和它的作用。
 */

// 项目中所有存储在localStorage中的数据
const localKeys = [
	// 主题色
	'theme',
	//  是否折叠侧边菜单
	'sideCollapse',
	// 是否显示顶部页面tag标签
	'tagVisible',
	// 系统风格 round,square
	'style',
	// 组件大小
	'size',
]

// 项目中所有存储在sessionStorage中的数据
const sessionKeys = ['token']

class StorageMng {
	// key名称前缀
	prefix = ''
	// 使用localStorage还是sessionStorage
	mode = localStorage

	constructor(mode, prefix = '') {
		this.prefix = prefix
		this.mode = mode
	}

	setItem(key, value) {
		try {
			this.mode.setItem(`${this.prefix}${key}`, window.JSON.stringify(value))
		} catch (err) {
			console.warn(`Storage ${key} set error`, err)
		}
	}

	getItem(key) {
		const result = this.mode.getItem(`${this.prefix}${key}`)
		try {
			return result ? window.JSON.parse(result) : result
		} catch (err) {
			console.warn(`Storage ${key} get error`, err)
		}
	}

	removeItem(key) {
		this.mode.removeItem(`${this.prefix}${key}`)
	}

	clear() {
		this.mode.clear()
	}

	getKey(index) {
		return this.getKeys()[index]
	}

	// 获取所有数据的名称
	getKeys() {
		const keys = []
		Array.from({ length: this.mode.length }).forEach((item, index) => {
			const key = this.mode.key(index)
			if (key.startsWith(this.prefix)) {
				keys.push(key.slice(this.prefix.length))
			}
		})
		return keys
	}

	// 获取所有数据
	getAll() {
		return Object.fromEntries(this.getKeys().map(key => [key, this.getItem(key)]))
	}
}

const localMng = new StorageMng(localStorage)
const sessionMng = new StorageMng(sessionStorage)

export { StorageMng, localMng, sessionMng }
