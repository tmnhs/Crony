<template>
	<el-radio-group v-model="systemStyle">
		<el-radio label="round">{{$t('setting.round')}}</el-radio>
		<el-radio label="square">{{$t('setting.square')}}</el-radio>
	</el-radio-group>
</template>

<script>
/**
 * 	 设置系统风格。按钮,表单可以设置成方形或椭圆形
 */

const roundContent = require('!css-loader!./style.css').default[0][1]

export default {
	data() {
		return {
			systemStyle: this.$store.getters.style,
		}
	},
	watch: {
		systemStyle(value) {
			this.triggerStyle(value)
			this.$store.commit('SET_STYLE', value)
		},
	},
	mounted() {
		this.triggerStyle(this.systemStyle)
	},
	methods: {
		triggerStyle(value) {
			const systemStyleElem = document.head.querySelector('#systemStyle')
			if (value === 'round') {
				if (systemStyleElem) return
				const styleElem = document.createElement('style')
				styleElem.setAttribute('id', 'systemStyle')
				styleElem.appendChild(document.createTextNode(roundContent))
				document.head.appendChild(styleElem)
			} else if (value === 'square') {
				if (systemStyleElem) {
					document.head.removeChild(systemStyleElem)
				}
			}
		},
	},
}
</script>
