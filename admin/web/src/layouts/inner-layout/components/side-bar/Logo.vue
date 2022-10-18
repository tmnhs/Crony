<template>
	<div class="side-bar-logo" :class="{ collapse: sideCollapse }">
		<router-link class="link" to="/">
			<base-icon name="logo" :size="24"></base-icon>
			<span v-show="!sideCollapse">{{$t('title')}}</span>
		</router-link>
	</div>
</template>

<script>
export default {
	computed: {
		sideCollapse() {
			return this.$store.getters.sideCollapse
		},
	},
}
</script>

<style lang="scss" scoped>
.side-bar-logo {
	&.collapse .link {
		width: 64px;
	}

	.link {
		display: block;
		width: 200px;
		height: 50px;
		line-height: 50px;
		background-color: #434a50;
		font-size: 18px;
		color: #fff;
		text-align: center;
		cursor: pointer;
		/* 查看element的menu.css源码可知transition的参数 */
		transition: width 0.3s ease-in-out;
	}
}
</style>
