<template>
	<div class="mine">
		<el-form :model="userInfo" :rules="rules" ref="userInfo" label-width="70px">
			<el-form-item :label="$t('user.name')" prop="username">
				<el-input v-model="userInfo.username"></el-input>
			</el-form-item>

			<el-form-item :label="$t('user.role')">
				<el-select v-model="userInfo.role"  readonly disabled >
					<el-option
						v-for="item in tableMng.getTable('role')"
						:key="item.id"
						:label="item.name"
						:value="item.id"
					></el-option>
				</el-select>
			</el-form-item>

			<!-- <el-form-item label="头像：">
				<avatar-upload v-model="userInfo.avatar" action="https://sm.ms/api/v2/upload" name="smfile" />
			</el-form-item> -->

			<el-row>
				<el-col :span="12">
					<el-form-item :label="$t('user.email')" :placeholder="$t('user.email_form')" prop="email">
						<el-input v-model="userInfo.email" clearable></el-input>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>

		<el-button type="primary" round :loading="submitLoading" @click="handleSubmit" > {{$t('common.submit')}} </el-button>
	</div>
</template>

<script>
import tableMng from '@/utils/tableMng'
import AvatarUpload from '@/components/business/upload/avatar-upload'

export default {
	name: 'Mine',
	components: {
		AvatarUpload,
	},
	data() {
		return {
			tableMng,
			rules: {
				name: [
							{
					required: true,
					message: 'please fill in name',
					trigger: 'blur',
				},
				{
					max: 20,
					message: 'the name must be no longer than 20 characters',
					trigger: 'blur',
				},
				],
			
				email: [
						{
						required: true,
						message: 'please fill in  email',
						trigger: 'blur',
					},
					{
						type: 'email',
						message: 'the email format is incorrect',
						trigger: 'blur',
					},
				],
			},
			submitLoading: false,
		}
	},
	computed: {
		userInfo() {
			return this.$store.getters.userInfo
		},
	},
	methods: {
		 handleSubmit() {
			this.submitLoading = true
			this.$refs.userInfo.validate(async valid => {
				if (valid) {
					const res= await this.$api.user.update( this.userInfo )
					this.$message.success('更新成功')
				} else {
				}
				this.submitLoading = false
			})
		},
	},
}
</script>
<style lang="scss" scoped>
.mine {
	width: 50%;
	min-width: 600px;
	margin: 0 auto;
	padding: 20px;
	background-color: #fff;
	border-radius: 10px;

	.el-button {
		display: block;
		margin: 0 auto;
	}
}
</style>
