<template>
	<div class="form-base">
		<el-form ref="form" :model="jobDetail" :rules="formRules" label-width="100px" class="job-form">

			<el-row :gutter="20">
				<el-col :xs="24" :md="12">
					<el-form-item :label="$t('job.name')" prop="name">
						<el-input v-model="jobDetail.name" :placeholder="$t('job.name_form')" clearable></el-input>
					</el-form-item>
				</el-col>
			</el-row>
			<el-row :gutter="20">
				<el-col :xs="24" :md="12">
					<el-form-item class="inherit-line-height" :label="$t('job.type')"  prop="job_type">
						<el-radio-group v-model.number="jobDetail.job_type">
							<el-radio :label="1">Shell Command</el-radio>
							<el-radio :label="2">Http CallBack</el-radio>
						</el-radio-group>
						 <el-tooltip class="item" effect="dark" :content="$t('job.type_tip')" placement="top-start">
							<i class="el-icon-warning-outline"></i>
						</el-tooltip>
					</el-form-item>
					
				</el-col>
				<el-col :xs="24" :md="12">
						<el-form-item class="inherit-line-height" :label="$t('job.http_method')" label-width="100px" prop="resource" v-if="jobDetail.job_type==2">
							<el-radio-group v-model.number="jobDetail.http_method">
								<el-radio :label="1">Get</el-radio>
								<el-radio :label="2">Post</el-radio>
							</el-radio-group>
						</el-form-item>
				</el-col>
			</el-row>
      <el-row :gutter="20" >
        <el-form-item :label="$t('job.script')"  prop="script" v-if="jobDetail.job_type==1">
          <el-select v-model="jobDetail.script_id" :placeholder="$t('job.script_form')"  multiple clearable>
            <el-option
                v-for="item in scriptList"
                :key="item.id"
                :label="item.name"
                :value="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
      </el-row>
				<el-form-item :label="$t('job.command')" prop="command">
				<el-input
					v-model="jobDetail.command"
					type="textarea"
					:rows="2"
					show-word-limit
					:placeholder="$t('job.command_form')" 
				></el-input>
			</el-form-item>
			<el-divider content-position="center">{{$t('job.cron_conf')}}</el-divider>
			<el-row class="cron-form">
				<el-col :xs="24" :md="12">
					<el-form-item :label="$t('job.cron')"  prop="spec" label-width="120px">
						<el-input v-model="jobDetail.spec" :placeholder="$t('job.cron_form')"  clearable></el-input>
					</el-form-item>
				</el-col>
				<el-col :xs="24" :md="12">
					<el-popover
						ref="ref_cron"
						placement="right"
						width="700"
						trigger="click">
						<Crontab
						@hide="handleCrontabHide"
						@fill="handleCrontabConfirm"
						>
						</Crontab>
						<el-tag class="cron-tag" slot="reference" type="success" >{{$t('job.cron_produce')}}</el-tag>
					</el-popover>
				</el-col>
			</el-row>
			
				
			<el-divider content-position="center">{{$t('job.allocation_conf')}}</el-divider>
			<el-row :gutter="20">
				<el-col :xs="24" :md="12">
						<el-form-item class="inherit-line-height" :label="$t('job.allocation')" label-width="100px" prop="allocation">
							<el-radio-group v-model.number="jobDetail.allocation">
								<el-radio :label="1">{{$t('job.manual')}}</el-radio>
								<el-radio :label="2">{{$t('job.auto')}}</el-radio>
					</el-radio-group>
					 <el-tooltip class="item" effect="dark" :content="$t('job.allocation_tip')" placement="top-start">
							<i class="el-icon-warning-outline"></i>
						</el-tooltip>
					</el-form-item>
				</el-col>
			</el-row>
		
			<el-row :gutter="20">
				<el-col :xs="24" :md="12">
					<el-form-item :label="$t('job.run_on')"  prop="run_on" v-if="jobDetail.allocation==1">
						<el-select v-model="jobDetail.run_on" :placeholder="$t('job.uuid_form')"  muliple clearable>
							<el-option
								v-for="item in nodeList"
								:key="item.uuid"
								:label="item.ip+':'+item.pid"
								:value="item.uuid"
							></el-option>
						</el-select>
					</el-form-item>
				</el-col>
			</el-row>
			<el-divider content-position="center">{{$t('job.job_conf')}}</el-divider>

			<el-row>
				<el-col :xs="24" :md="8">
					<el-form-item :label="$t('job.timeout')"  prop="timeout" >
						<el-input v-model.number="jobDetail.timeout"  clearable></el-input>
									</el-form-item>
					
				</el-col>
				<el-col :xs="24" :md="8">
						<el-form-item>
							<el-tooltip effect="dark" :content="$t('job.job_tip')" placement="top-start">
							<i class="el-icon-warning-outline"></i>
						</el-tooltip>
						</el-form-item>
				</el-col>
			</el-row>
			<el-row :gutter="20">
				
				<el-col :xs="24" :md="12">
					<el-form-item :label="$t('job.retry_time')"  prop="retry_times" >
						<el-input v-model.number="jobDetail.retry_times"  clearable></el-input>
					</el-form-item>
				</el-col>
					<el-col :xs="24" :md="12">
					<el-form-item :label="$t('job.retry_interval')"  prop="retry_interval" label-width="110px" >
						<el-input v-model.number="jobDetail.retry_interval" maxlength="120px" placeholder="设置任务重试时间间隔" clearable></el-input>
					</el-form-item>
				</el-col>

			</el-row>
			<el-form-item :label="$t('job.note')" prop="note">
				<el-input
					v-model="jobDetail.note"
					type="textarea"
					:rows="1"
					show-word-limit
				></el-input>
			</el-form-item>
			<el-divider content-position="center">{{$t('job.notify_conf')}}</el-divider>
			<el-row :gutter="20">
				<el-col :xs="24" :md="14">
						<el-form-item class="inherit-line-height" :label="$t('job.notify_method')" label-width="120px"   prop="notify_type">
							<el-radio-group v-model="jobDetail.notify_type">
								<el-radio :label="1">Email</el-radio>
								<el-radio :label="2">Webhook</el-radio>
					</el-radio-group>
					<el-tooltip effect="dark" :content="$t('job.notify_tip')" placement="top-start">
							<i class="el-icon-warning-outline"></i>
					</el-tooltip>
					</el-form-item>
				</el-col>
			</el-row>
			<el-row :gutter="20" >	
					<el-form-item :label="$t('job.notify_to')"  prop="notify_to">
						<el-select v-model="jobDetail.notify_to" :placeholder="$t('job.notify_form')"  multiple clearable>
							<el-option
								v-for="item in userList"
								:key="item.id"
								:label="item.username"
								:value="item.id"
							></el-option>
						</el-select>
			</el-form-item>
			</el-row>
		</el-form>

		<div style="text-align: center">
			<el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{$t('common.submit')}}</el-button>
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
	job_type: null,
	command: '',
	run_on:'',
	timeout:0,
	retry_times:3,
	retry_interval:1,
	notify_type:1,
  notify_to:[],
  script_id:[],
	spec:'',
	allocation:1,
	note:'',
}
export default {
	name: 'JobEdit',
	props: ['job_id'],
	components:{
		Crontab
	},
	data() {
		return {
			tableMng,
			jobDetail: { ...defaultDetail },
      userList:[],
      scriptList:[],
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
		
				job_type: [
					{
						required: true,
						message: 'Please select a job type',
						trigger: 'blur',
					},
				],
				command: [
					{
						required: true,
						message: 'Please enter the job command',
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
		this.getUserList()
		this.getNodeList()
    this.getScriptList()
	},
	methods: {
		async getUserList() {
			const data = await this.$api.user.getList({page:1,page_size:1000})
			this.userList = data.list
		},
    async getScriptList() {
      const data = await this.$api.script.getScriptList({page:1,page_size:1000})
      this.scriptList = data.list
    },
		//获取节点列表
		async getNodeList() {
			const data = await this.$api.node.getNodeList({status:1})
			this.nodeList = data.list
		},
		async getDetail() {
			if (this.job_id) {
				const data = await this.$api.job.findJob({ id: this.job_id })
				this.jobDetail = {
					id: data.id,
					name: data.name,
					job_type: data.job_type,
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
          script_id: data.script_id,
				}
			} else {
				this.jobDetail = { ...defaultDetail }
			}
		},
		handleSubmit() {
			this.$refs.form.validate(async valid => {
				if (valid) {
					this.confirmLoading = true
					try {
						await this.$api.job.updateJob(this.jobDetail )
						this.$message.success('create/update success')
						this.$emit('success')
						this.$router.push('/job/list')
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
			this.jobDetail.spec=spec
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
.job-form label{
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
