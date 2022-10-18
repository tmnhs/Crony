<template>
	<div class="login">
		<p class="login__header">{{$t('login.title')}}</p>

		<el-form
			class="login__body"
			ref="formRef"
			:model="formData"
			:rules="formRules"
			auto-complete="on"
			label-width="80px"
			label-position="top"
		>
			<el-form-item :label="$t('login.account')" prop="username">
				<el-input type="text" auto-complete="on" autofocus v-model="formData.username" :placeholder="$t('login.account_form')">
					<base-icon name="user" slot="prefix" icon-class="icon"></base-icon>
				</el-input>
			</el-form-item>

			<el-form-item :label="$t('user.password')" prop="password">
				<el-input
					:type="passwordType"
					auto-complete="on"
					v-model="formData.password"
					@keyup.enter.native="handleLogin"
					:placeholder="$t('user.password_form')"
				>
					<base-icon slot="prefix" name="lock" icon-class="icon"></base-icon>
					<base-icon
						slot="suffix"
						:name="passwordType === 'password' ? 'eye-close' : 'eye-open'"
						icon-class="icon icon-eye"
						@click.native="showPwd"
					></base-icon>
				</el-input>
			</el-form-item>

			<el-form-item>
				<el-checkbox v-model="formData.rememberPwd">{{$t('login.remember')}}</el-checkbox>
				
			</el-form-item>

			<el-form-item>
				<el-button type="primary" :loading="loginLoading" @click="handleLogin">{{$t('login.login')}}</el-button>
			</el-form-item>
		</el-form>
	</div>
</template>

<script>
import { encryptByDES } from '@/utils/crypto'

export default {
	data() {
		return {
			formData: {
				username: 'root',
				password: '123456',
				rememberPwd: false,
			},
			formRules: {
				username: [
					{
						required: true,
						message: '账号不能为空',
						trigger: 'blur',
					},
				],
				password: [
					{
						required: true,
						message: '请输入密码',
						trigger: 'blur',
					},
					{
						min: 6,
						message: '密码长度不能少于六位',
						trigger: 'blur',
					},
				],
			},
			passwordType: 'password',
			loginLoading: false,
		}
	},
	methods: {
		showPwd() {
			this.passwordType = this.passwordType == 'password' ? 'text' : 'password'
		},
		handleLogin() {
			this.$refs.formRef.validate(async valid => {
				if (valid) {
					try {
						this.loginLoading = true
						const account = {
							username: this.formData.username,
							password: this.formData.password,
						}
						await this.$store.dispatch('login', account)
					} catch (err) {
						console.error(err)
					} finally {
						this.loginLoading = false
					}
				}
			})
		},
	},
}
</script>

<style lang="scss" scoped>
.login {
	padding: 15px 20px;
	background-color: rgba(255, 255, 255, 0);
	box-shadow: 0px 0px 20px 0px rgba(0, 0, 0, 0.7);
	border-radius: 10px;

	.login__header {
		font-size: 22px;
		text-align: center;
		color: #fff;
	}

	.login__body {
		.icon {
			color: $black3;
			font-size: 18px;
		}

		.icon.icon-eye {
			padding-right: 10px;
			cursor: pointer;
		}

		.register {
			float: right;
			color: #fff;
		}
	}
}
</style>

<style lang="scss">
.login {
	.el-input--prefix {
		.el-input__prefix {
			left: 10px;
		}

		.el-input__inner {
			padding-left: 35px;
			padding-right: 35px;
		}
	}

	.el-form-item__label,
	.el-checkbox__label {
		color: #fff;
	}

	.el-button {
		display: block;
		width: 100%;
	}
}
</style>
