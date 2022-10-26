<template>
	<div class="form-base">
		<el-form ref="form" :model="scriptDetail" :rules="formRules" label-width="100px" class="script-form">

			<el-row :gutter="20">
				<el-col :xs="24" :md="12">
					<el-form-item :label="$t('script.name')" prop="name">
						<el-input v-model="scriptDetail.name" :placeholder="$t('script.name_form')" clearable></el-input>
					</el-form-item>
				</el-col>
			</el-row>
				<el-form-item :label="$t('script.command')" prop="command">
				<el-input
					v-model="scriptDetail.command"
					type="textarea"
					:rows="8"
					show-word-limit
					:placeholder="$t('script.command_form')" 
				></el-input>
			</el-form-item>


		</el-form>

		<div style="text-align: center">
			<el-button type="primary" @click="handleSubmit" :loading="submitLoading" >{{$t('common.submit')}}</el-button>
			<el-button type="info" @click="handleCancel">{{$t('common.cancel')}}</el-button>
		</div>
	</div>
</template>

<script>
import tableMng from '@/utils/tableMng'
import Crontab from '@/components/crontab/index.vue'
const defaultDetail = {
	id: null,
	name: '',
	script_type: null,
	command: '',
	run_on:'',
	timeout:0,
	retry_times:3,
	retry_interval:1,
	notify_type:1,
	notify_to:[],
	spec:'',
	allocation:1,
	note:'',
}
export default {
	name: 'scriptEdit',
	props: ['script_id'],
	components:{
		Crontab
	},
	data() {
		return {
			tableMng,
			scriptDetail: { ...defaultDetail },
			userList:[],
			nodeList:[],
			cronTable:false,
			formRules: {
				name: [
					{
						required: true,
						message: 'Please fill in name',
						trigger: 'blur',
					},
					{
						max: 20,
						message: 'The name cannot exceed 20 characters',
						trigger: 'blur',
					},
				],
		
				script_type: [
					{
						required: true,
						message: 'Please select a script type',
						trigger: 'blur',
					},
				],
				command: [
					{
						required: true,
						message: 'Please enter the script command',
						trigger: 'blur',
					},
				],
				spec: [
				{
					required: true,
					message: 'Please enter cron expression',
					trigger: 'blur',
				},
			],
				notify_type: [
					{
						required: true,
						message: 'Please select notify type',
						trigger: 'blur',
					},
				],
				run_on: [
				{
					required: true,
					message: 'Please select a running node',
					trigger: 'blur',
				},
				]
			},
			regionList: [],
			submitLoading: false,
			visible: false,
			confirmLoading: false,
		}
	},
	created() {
		this.getDetail()
		this.getNodeList()
	},
	methods: {
		//获取节点列表
		async getNodeList() {
			const data = await this.$api.node.getNodeList({status:1})
			this.nodeList = data.list
		},
		async getDetail() {
			if (this.script_id) {
				const data = await this.$api.script.findScript({ id: this.script_id })
				this.scriptDetail = {
					id: data.id,
					name: data.name,
					script_type: data.script_type,
					http_method:data.http_method,
					command: data.command,
					run_on:data.run_on,
					timeout:data.timeout,
					retry_times:data.retry_times,
					retry_interval:data.retry_interval,
					notify_to:data.notify_to,
					notify_type:data.notify_type,
					spec:data.spec,
					allocation:1,
					note:data.note,
				}
			} else {
				this.scriptDetail = { ...defaultDetail }
			}
		},
		handleSubmit() {
			this.$refs.form.validate(async valid => {
				if (valid) {
					this.confirmLoading = true
					try {
						await this.$api.script.updateScript(this.scriptDetail )
						this.$message.success('create/update success')
						this.$emit('success')
						this.$router.push('/script/list')
					} catch (err) {
						console.error(err)
					} finally {
						this.confirmLoading = false
					}
				} else {
					this.$message.error($t('common.role_error'))
				}
			})
		},
		handleCancel() {
			this.$refs.form.resetFields()
		},

		handleCrontabConfirm(spec) {
			this.scriptDetail.spec=spec
			this.$refs.ref_cron.doClose()
		},
		handleCrontabHide(){
			this.$refs.ref_cron.doClose()
		},
	},
}
</script>

<style lang="scss" scoped>
.form-base {
	width: 70%;
	min-width: 500px;
	padding: 1em;
	margin: 0 auto;
	border-radius: 10px;
	background-color: #fff;
}
.el-row {
	margin-bottom: 20px;
}
</style>

<style lang="scss">
.form-base {
	.inherit-line-height {
		.el-form-item__label,
		.el-form-item__content {
			line-height: inherit;
		}
	}
}
.script-form label{
    color: #99a9bf;
}
.cron-tag{
	width: 130px;
	margin-left: 15px;
	margin-bottom: 20px;
	height: 30px;
	line-height: 30px;
	float: right;
}

</style>
