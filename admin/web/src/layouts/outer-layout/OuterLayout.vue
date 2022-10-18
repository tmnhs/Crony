<template>
	<div class="outer-layout">
		<div class="outer-layout__header">
			<!-- todo icon -->
			<base-icon name="logo" :size="36"></base-icon>
			<span>{{$t('title')}}</span>
		</div>

		<div class="outer-layout__page">
			<router-view></router-view>
		</div>
	</div>
</template>

<script>
/**
 * 外部布局(未登录)
 */

export default {
	data() {
		return {}
	},
}
</script>

<style lang="scss" scoped>
.outer-layout {
	min-height: 100vh;
	background: url(~@/assets/images/layout/outer_bg.jpg) no-repeat fixed center;
	background-size: cover;

	&__header {
		padding-top: 60px;
		text-align: center;
		font-size: 32px;
		color: #fff;
	}

	&__page {
		width: 440px;
		margin: 0 auto;
		margin-top: 60px;
	}
}
</style>
