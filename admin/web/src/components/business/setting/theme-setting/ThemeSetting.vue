<template>
	<el-color-picker popper-class="theme-setting" v-model="theme" :predefine="predefineThemes" size="small" />
</template>

<script>
/**
 * 换肤
 */

// 获取element的样式
const elementStyleContent = require('!css-loader!element-ui/lib/theme-chalk/index.css').default[0][1]

export default {
	data() {
		return {
			theme: '',
			predefineThemes: [
				'#409eff',
				'#1cc09d',
				'#ffa69e',
				'#ff4879',
				'#373737',
				'#5b5bea',
				'#FF8822',
				'#757575',
				'#5FC3D7',
				'#ffd166',
			],
			defaultStyle: elementStyleContent,
			defaultColors: [],
		}
	},
	watch: {
		theme(newTheme) {
			this.updateElementTheme(newTheme)
			this.updateCustomTheme(newTheme)
			this.$store.commit('SET_THEME', newTheme)
		},
	},
	created() {
		this.defaultColors = this.getColors('409eff')
		// 生成的颜色系列：["409eff", "64,158,255", "53a8ff", "66b1ff", "79bbff", "8cc5ff", "a0cfff", "b3d8ff", "c6e2ff", "d9ecff", "ecf5ff", "3a8ee6"]
		this.theme = this.$store.getters.theme
	},
	mounted() {
		this.createStyleElement('elementTheme')
	},
	methods: {
		createStyleElement(id) {
			if (document.getElementById(id)) return
			const styleElem = document.createElement('style')
			styleElem.setAttribute('id', id)
			document.head.appendChild(styleElem)
		},
		// 由基础颜色值生成一系列颜色值
		getColors(theme) {
			// 实现scss的mix函数，主题色与#fff进行混合。与白色混合其实就是改变透明度。
			// mix(#fff, #409eff, 90%)和 rgba(64,158,255,0.1); 它们的效果是一样的

			const tintColor = (color, tint) => {
				let red = parseInt(color.slice(0, 2), 16)
				let green = parseInt(color.slice(2, 4), 16)
				let blue = parseInt(color.slice(4, 6), 16)
				// 有些背景需要设置透明度，用到了rgba颜色值。
				if (tint === 0) {
					return [red, green, blue].join(',')
					//如果是未经压缩的css文件，或scss编译后的css文件，注意rgba中逗号后会有一个空格。
				} else {
					red = Math.round(red * (1 - tint) + 255 * tint)
					green = Math.round(green * (1 - tint) + 255 * tint)
					blue = Math.round(blue * (1 - tint) + 255 * tint)
					red = red.toString(16)
					green = green.toString(16)
					blue = blue.toString(16)
					return `${red}${green}${blue}`
				}
			}
			// $button-active使用的这个色
			const shadeColor = (color, shade) => {
				let red = parseInt(color.slice(0, 2), 16)
				let green = parseInt(color.slice(2, 4), 16)
				let blue = parseInt(color.slice(4, 6), 16)
				red = Math.round((1 - shade) * red)
				green = Math.round((1 - shade) * green)
				blue = Math.round((1 - shade) * blue)
				red = red.toString(16)
				green = green.toString(16)
				blue = blue.toString(16)
				return `${red}${green}${blue}`
			}

			const colors = [theme]
			for (let i = 0; i < 10; i++) {
				colors.push(tintColor(theme, i / 10))
			}
			colors.push(shadeColor(theme, 0.1))
			return colors
		},
		// 更新element主题
		updateElementTheme(newTheme) {
			const newColors = this.getColors(newTheme.replace('#', ''))
			let newStyle = this.defaultStyle
			this.defaultColors.forEach((color, index) => {
				newStyle = newStyle.replace(new RegExp(color, 'ig'), newColors[index])
			})
			document.head.querySelector('#elementTheme').innerText = newStyle
		},
		// 更新自己书写的css的主题
		updateCustomTheme(newTheme) {
			const newColors = this.getColors(newTheme.replace('#', ''))
			const rootStyle = document.documentElement.style
			newColors.forEach((color, index) => {
				if (index === 0) {
					rootStyle.setProperty('--theme', `#${color}`)
					return
				}
				if (index === 11) {
					rootStyle.setProperty('--theme-shade', `#${color}`)
					return
				}
				if (index >= 2) {
					rootStyle.setProperty(`--theme-white__${index - 1}`, `#${color}`)
				}
			})
		},
	},
}
</script>

<style lang="scss">
.theme-setting {
	.el-color-dropdown__value,
	.el-color-dropdown__link-btn {
		display: none;
	}
}
</style>
