<template>
	<el-select v-model="size" :key="sizeSelectKey">
		<el-option label="big" value="default"></el-option>
		<el-option label="medium" value="medium"> </el-option>
		<el-option label="small" value="small"></el-option>
		<el-option label="mini" value="mini"></el-option>
	</el-select>
</template>

<script>
/**
 * ElementUI元素大小设置
 */
export default {
	data() {
		return {
			size: this.$store.getters.size,
			sizeSelectKey: 0,
		}
	},
	watch: {
		size(value) {
			this.$ELEMENT.size = value
			this.sizeSelectKey++
			this.$store.commit('SET_SIZE', value)
			this.$store.commit('SET_CACHE_PAGES', [])
			// this.$router.replace('/reload')
		},
	},
}
</script>
