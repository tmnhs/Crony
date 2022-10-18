<template>
	<base-drawer :title="$t('setting.title')" :visible="visible" size="300px" :show-footer="false" @close="visible = false">
		<el-form label-width="100px">
			<el-form-item :label="$t('setting.language')">
						<el-select v-model="language" :key="languageSelectKey">
							<el-option label="中文" value="zh"></el-option>
							<el-option label="English" value="en"> </el-option>
						</el-select>
			</el-form-item>
			<el-form-item :label="$t('setting.navigation')">
				<el-switch v-model="tagVisible"></el-switch>
			</el-form-item>

			<el-form-item :label="$t('setting.style')">
				<style-setting></style-setting>
			</el-form-item>

			<el-form-item :label="$t('setting.size')">
				<size-setting></size-setting>
			</el-form-item>

			<el-form-item :label="$t('setting.skin')">
				<theme-setting></theme-setting>
			</el-form-item>

		</el-form>
	</base-drawer>
</template>

<script>
import StyleSetting from '@/components/business/setting/style-setting'
import SizeSetting from '@/components/business/setting/size-setting'
import ThemeSetting from '@/components/business/setting/theme-setting'

export default {
	components: {
		StyleSetting,
		SizeSetting,
		ThemeSetting,
	},
	data() {
		return {
			visible: false,
			language:this.$store.getters.language,
			tagVisible: this.$store.getters.tagVisible,
			languageSelectKey:0,
		}
	},
	watch: {
		tagVisible(value) {
			this.$store.commit('SET_TAG_VISIBLE', value)
		},
		language(value) {
			console.log(value);
			this.language = value
			this.languageSelectKey++
			this.$i18n.locale = value
			this.$store.dispatch('setLanguage', value)
			// this.$router.replace('/reload')
		},
	},
	methods: {
		open() {
			this.visible = true
		},
	},
}
</script>
