<template>
	<div>
		<router-view></router-view>

		<div style="display: none">
			<style-setting></style-setting>
			<theme-setting></theme-setting>
		</div>
	</div>
</template>

<script>
import StyleSetting from '@/components/business/setting/style-setting'
import ThemeSetting from '@/components/business/setting/theme-setting'

export default {
	components: {
		StyleSetting,
		ThemeSetting,
	},
}
</script>
