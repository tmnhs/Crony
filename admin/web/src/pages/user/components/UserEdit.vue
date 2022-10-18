<template>
	<base-dialog
		class="user-edit"
		:title="detail.id ? 'Edit User' : 'Add User'"
		width="700px"
		:visible="visible"
		:confirm-loading="confirmLoading"
		@confirm="handleConfirm"
		@close="handleBeforeClose"
	>
		<el-form ref="form" :model="detail" :rules="formRules" label-width="90px" v-loading="initLoading">
			<el-form-item :label="$t('user.name')" prop="username">
				<el-input v-model="detail.username" :placeholder="$t('user.name_form')" clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('user.password')"  prop="password" v-if="detail.id?false:true">
				<el-input v-model="detail.password"  :placeholder="$t('user.password_form')" clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('user.role')" prop="role">
				<el-select v-model="detail.role"  :placeholder="$t('user.role_query')">
					<el-option
						v-for="item in tableMng.getTable('role')"
						:key="item.id"
						:label="item.name"
						:value="item.id"
					></el-option>
				</el-select>
			</el-form-item>

			<el-row>
				<el-col :span="12">
					<el-form-item :label="$t('user.email')"  prop="email">
						<el-input v-model="detail.email" :placeholder="$t('user.email_form')"  clearable></el-input>
					</el-form-item>
				</el-col>
			</el-row>
		</el-form>
	</base-dialog>
</template>

<script>
import _ from 'lodash'
import tableMng from '@/utils/tableMng'
import i18n from  "@/assets/lang";
const defaultDetail = {
	id: null,
	username: '',
	role: null,
	email: '',
}

const formRules = {
	username: [
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
	password: [
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
	role: [
		{
			required: true,
			message: 'please select roles',
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
}

export default {
	name:"UserEdit",
	data() {
		return {
			tableMng,
			formRules,
			visible: false,
			detail: {},
			detailBackup: {},
			confirmLoading: false,
			initLoading: false,
		}
	},
	methods: {
		// 打开模态窗
		open(data = {}) {
			this.visible = true
			this.changeDetail(data)
			this.getDetail()
		},
		toggleVisible() {
			this.visible = !this.visible
		},
		changeDetail(data) {
			this.detail = { ...defaultDetail, ...data }
			this.detailBackup = { ...defaultDetail, ...data }
		},
		// 获取用户详情
		async getDetail() {
			const { id } = this.detail
			if (!id) return
			this.initLoading = true
			try {
				const res = await this.$api.user.getDetail({ id })
				this.changeDetail(res)
			} catch (err) {
				console.error(err)
			} finally {
				this.initLoading = false
			}
		},
		handleConfirm() {
			this.$refs.form.validate(async valid => {
				if (valid) {
				if (this.detail.id==0||this.detail.id==null){
					await this.$api.account.register(this.detail )
					this.$message.success(i18n.t('common.add_success'))
				}else {
					await this.$api.user.update( this.detail )
					this.$message.success(i18n.t('common.edit_success'))
				}
				this.handleClose()
				this.$emit('success')
		
				} else {
					this.$message.error(i18n.t('common.role_error'))
				}
			})
		},
		handleBeforeClose() {
			if (window.JSON.stringify(this.detailBackup) === window.JSON.stringify(this.detail)) {
				this.handleClose()
			} else {
				this.handleClose()
			}
		},
		handleClose() {
			this.$refs.form.clearValidate()
			this.visible = false
		},
	},
}
</script>
