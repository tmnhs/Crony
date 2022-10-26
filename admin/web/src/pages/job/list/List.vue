<template>
	<div class="job-list" >
		<div class="job-list__header">
			<div class="title">
				<section-title :name="$t('job.job_list')" />
			</div>

			<div class="operation">
				<el-button type="primary" icon="el-icon-plus" @click="handleAdd">{{$t('common.add')}}</el-button>
				<el-button type="danger" icon="el-icon-minus" @click="handleDelete(null)" >{{$t('common.deletes')}}</el-button>
				<el-button icon="el-icon-download" :loading="exportLoading" @click="handleExport">{{$t('common.export')}}</el-button>

			</div>
		</div>

		<el-form :inline="true" :model="query">
			<el-form-item :label="$t('job.id')">
				<el-input v-model.number="query.id" :placeholder="$t('job.id_form')" clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('job.name')">
				<el-input v-model="query.name" :placeholder="$t('job.name_form')" clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('job.uuid')">
				<el-input v-model="query.run_on" :placeholder="$t('job.uuid_form')"  clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('job.type')">
				<el-select v-model="query.job_type" placeholder="Select"  clearable>
					<el-option
						v-for="item in tableMng.getTable('job_type')"
						:key="item.id"
						:label="item.name"
						:value="item.id"
					></el-option>
				</el-select>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-search" @click="getjobList"> {{$t('common.search')}} </el-button>
				<el-button icon="el-icon-search" @click="handleReset">{{$t('common.reset')}}</el-button>

			</el-form-item>
		</el-form>

		<div class="job-list__table">
			<el-table
				:data="jobList"
				v-loading="jobTableLoading"
				@selection-change="handleSelectedRows"
			>	
				<el-table-column type="selection" width="50"></el-table-column>
				<el-table-column type="expand">
					<template slot-scope="scope">
						<el-form label-position="left" inline class="job-table-expand">
							<el-row :gutter="20">
								<el-col :xs="24" :md="12">
									<el-form-item :label="$t('job.id')">
										<span>{{ scope.row.id }}</span>
									</el-form-item>
								</el-col>
								<el-col :xs="24" :md="12">
													
									<el-form-item :label="$t('job.name')">
										<span>{{ scope.row.name }}</span>
									</el-form-item>
								</el-col>
						</el-row>
						<el-row :gutter="20">
								<el-col :xs="24" :md="24">
									<el-form-item :label="$t('job.command')">
										<textarea  class="command-text" v-model="scope.row.command"></textarea> 
									</el-form-item>
								</el-col>
						</el-row>
						<el-row :gutter="20">
								<el-col :xs="24" :md="24">
									<el-form-item :label="$t('job.run_on')">
										<span>{{ scope.row.run_on }}</span>
									</el-form-item>
								</el-col>
						</el-row>
							<el-row :gutter="20">
								<el-col :xs="24" :md="12">
									<el-form-item :label="$t('job.cron')">
										<span>{{ scope.row.spec }}</span>
									</el-form-item>
								</el-col>
								<el-col :xs="24" :md="12">
													
									<el-form-item :label="$t('job.type')">
										<span>{{ tableMng.getNameById('job_type', scope.row.job_type) }}{{scope.row.job_type==2?tableMng.getNameById('http_method', scope.row.http_method):''}}</span>
									</el-form-item>
								</el-col>
						</el-row>
						<el-row :gutter="20">
								<el-col :xs="24" :md="12">
									<el-form-item :label="$t('job.timeout')">
										<span>{{scope.row.timeout}}</span>
									</el-form-item>
								</el-col>
								<el-col :xs="24" :md="12">			
									<el-form-item :label="$t('job.retry_time')">
										<span>{{ scope.row.retry_times }}</span>
									</el-form-item>
								</el-col>
						</el-row>
						<el-row :gutter="20">
								<el-col :xs="24" :md="12">
									<el-form-item :label="$t('job.retry_interval')">
										<span>{{ scope.row.retry_interval }}</span>
									</el-form-item>
								</el-col>
								<!-- <el-col :xs="24" :md="12">			
									<el-form-item label="重试次数">
										<span>{{ scope.row.retry_times }}</span>
									</el-form-item>
								</el-col> -->
						</el-row>
						<el-row :gutter="20">
								<el-col :xs="24" :md="12">
									<el-form-item :label="$t('job.notify_to')">
										<span>{{ scope.row.notify_to!=null ? scope.row.notify_to:"NULL" }}</span>
									</el-form-item>
								</el-col>
              <el-col :xs="24" :md="12">
                <el-form-item :label="$t('job.script')">
                  <span>{{ scope.row.script_id!=null ? scope.row.script_id:"NULL" }}</span>
                </el-form-item>
              </el-col>
						</el-row>
							<el-row :gutter="20">
								<el-col :xs="24" :md="12">
									<el-form-item :label="$t('job.note')">
										<span>{{scope.row.note }}</span>
									</el-form-item>
								</el-col>
								<el-col :xs="24" :md="12">
									<el-form-item :label="$t('job.created')">
										<span>{{formatDate( scope.row.created )}}</span>
									</el-form-item>
								</el-col>
			
						</el-row>
						</el-form>
					</template>
				</el-table-column>
				<el-table-column prop="id" :label="$t('job.id')" width="80px"></el-table-column>
				<el-table-column prop="name" :label="$t('job.name')" width="120px" :show-overflow-tooltip="true"></el-table-column>
				<el-table-column prop="command" :label="$t('job.command')" :show-overflow-tooltip="true"></el-table-column>
				<el-table-column prop="run_on" :label="$t('job.run_on')" width="350px"></el-table-column>
				<el-table-column prop="type" :label="$t('job.type')"  width="150px">
					<template slot-scope="scope">
						{{ tableMng.getNameById('job_type', scope.row.job_type) }}{{scope.row.job_type==2?tableMng.getNameById('http_method', scope.row.http_method):''}}
					</template>
				</el-table-column>
				<el-table-column label="操作" width="350px">
					<template slot-scope="scope">
						<router-link  :to="`/job/log/${scope.row.id}`">
							<el-button type="text"> {{$t('job.log')}} </el-button>
						</router-link>
						<el-divider direction="vertical"></el-divider>

						<el-button type="text" @click="onceExecute(scope.row)"> {{$t('job.once')}} </el-button>
						<el-divider direction="vertical"></el-divider>
						<router-link :to="`/job/edit/${scope.row.id}`">
							<el-button type="text">{{$t('common.edit')}}</el-button>
						</router-link>
					
						<el-divider direction="vertical"></el-divider>
						<el-button type="text" @click="handleDelete(scope.row)" > {{$t('common.delete')}} </el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>

		<pagination
			:total="total"
			:page-number.sync="query.page"
			:page-size.sync="query.page_size"
			@pagination="getjobList"
		></pagination>
	</div>
</template>

<script>
/**
 * 任务管理
 */
import { scroll } from '@/utils/core'
import tableMng from '@/utils/tableMng'
import { exportExcel } from '@/utils/excle'
import { formatDate } from '@/utils/date'
import i18n from  "@/assets/lang";

const defaultQuery = {
	id:null,
	name: '',
	job_type: null ,
	run_on:'',
	kind:null,
	page: 1,
	page_size: 10,
}
export default {
	name: 'jobList',
	props: ['uuid'],
	data() {
		return {
			tableMng,
			userLoading: false,
			jobList: [],
			jobTableLoading: false,
			query: _.cloneDeep(defaultQuery),
			total: 0,
			selectedRows: [],
			exportLoading: false,
		}
	},
	created() {
		this.getjobList()
	},
	methods: {
		formatDate(time) {
			time = time * 1000
			let date = new Date(time)
			return formatDate(date, 'yyyy-MM-dd hh:mm:ss')
      },
		// 获取任务列表
		async getjobList() {
			if(this.uuid){
				this.query.run_on=this.uuid
			}
      if(this.query.id==''){
        this.query.id=null
      }
			this.jobTableLoading = true
			const data = await this.$api.job.getJobList(this.query)
			this.jobList = data.list
			this.total = data.total
			this.jobTableLoading = false
			const scrollElement = document.querySelector('.inner-layout__page')
			scroll(scrollElement, 0, 800)
		},
		// 跳转到新建任务页面
		handleAdd() {
			this.$router.push('/job/add')
		},
		// 删除
		handleDelete(row) {
			let ids = []
			let name = []
			if (row) {
				ids = [row.id]
				name = [row.name]
			} else {
				ids = this.selectedRows.map(row => row.id)
				name = this.selectedRows.map(row => row.name)
			}
			if (ids.length === 0) {
				this.$message.warning(i18n.t('common.delete_warn'))
			} else {
				this.$confirm(i18n.t('common.delete_question')+`"${name.join(',')}”？`,  i18n.t('common.prompt'), {
					type: 'warning',
				})
					.then(async () => {
						await this.$api.job.removeJob({ ids:ids })
						this.$message.success(i18n.t('common.delete_success'))
						this.getjobList()
					})
					.catch(() => {})
			}
		},
		// 多选
		handleSelectedRows(rows) {
			this.selectedRows = rows
		},

		// 导出任务表格
		async handleExport() {
			this.exportLoading = true
			const header = ['JobId', 'Job名','Command','RunOn','Spec"','JobType','HttpMethod','Timeout','RetryTime','RetryInterval','NotifyType','NotifyTo', 'CratedTime']
			const res = await this.$api.job.getJobList()
			const data = res.list.map((item, index) => {
				return {
						id: item.id,
						name: item.name,
						command: item.command,
						run_on:item.run_on,
						spec:item.spec,
						job_type: item.job_type,
						http_method: item.http_method,
						timeout:item.timeout,
						retry_times:item.retry_times,
						retry_interval:item.retry_interval,
						notify_type:item.notify_type,
						notify_to:item.notify_to,
						created:item.created,
				}
			})
			exportExcel(header, data, 'Job Data Sheet')
			this.exportLoading = false
		},
		
		// 重置查询
		handleReset() {
			this.query = _.cloneDeep(defaultQuery)
			this.uuid=''
			this.getjobList()
		},
		//立即执行
		async  onceExecute(row){
			await this.$api.job.once({ job_id:row.id ,node_uuid:row.run_on })
			this.$router.push('/job/log/'+row.id)
		}
	},
}
</script>

<style lang="scss" >
.job-list {
	min-height: 100%;
	padding: 1em;
	box-sizing: border-box;
	background-color: #fff;

	.job-list__header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		flex-wrap: wrap;

		.title,
		.operation {
			margin-bottom: 1em;
		}
	}
}
.job-table-expand label{
	font-size: 12px;
	width: 90px;
    color: #99a9bf;
}
.job-table-expand {
	
	 .el-form-item{
		margin-right: 0;
		margin-bottom: 0;
		width: 50%;
	 }
	 span{
		font-size: 5px;
	 }
	 
}
.command-text{
	font-size: 10px;
	width: 450px;
	color: #fff;
	background-color: rgb(14, 13, 13);
	resize:none;
	border-radius: 5px;
}

</style>
