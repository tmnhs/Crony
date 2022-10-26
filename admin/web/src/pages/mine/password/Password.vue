<template>
	<div class="mine">
		<el-form :model="formData" ref="formData" :rules="rules" label-width="120px" >
			<el-form-item :label="$t('user.old_password')"  prop="password">
				<el-input v-model="formData.password" :placeholder="$t('user.password_form')" ></el-input>
			</el-form-item>

			<el-form-item :label="$t('user.new_password')" prop="new_password">
				<el-input v-model="formData.new_password" :placeholder="$t('user.password_form')" ></el-input>
			</el-form-item>
		</el-form>
		<el-button type="primary" style="margin-top:50px;" round :loading="submitLoading" @click="handleSubmit" disabled> {{$t('common.submit')}} </el-button>
	</div>
</template>

<script>

export default {
	name: 'Password',
	data() {
		return {
			formData:{
				password:'',
				new_password:'',
			},
			rules: {
					new_password: [
					{
						required: true,
						message: 'Please enter password',
						trigger: 'blur',
					},
					{
						min: 6,
						message: 'the password must contain at least six characters',
						trigger: 'blur',
					},
				],
				
			},
			submitLoading: false,
		}
	},
	methods: {
		 async handleSubmit() {
			// this.submitLoading = true
			this.$refs.formData.validate(async valid => {
				if (valid) {
					const res= await this.$api.account.changePassword( this.formData )
					this.$message.success('修改成功')
				} else {
				}
				// this.submitLoading = false
			})
		},
	},
}
</script>
<style lang="scss" scoped>
.mine {
	width: 40%;
	min-width: 400px;
	min-height: 180px;
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
