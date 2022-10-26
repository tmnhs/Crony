<template>
	<div class="log-manager">
		<div class="log-manager__header">
			<div class="title">
				<section-title :name="$t('log.log_list')" />
			</div>
		</div>

		<el-form ref="queryForm" :inline="true" :model="query">
			<el-form-item :label="$t('log.id')">
				<el-input v-model.number="query.job_id" :placeholder="$t('job.id_form')" clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('job.name')">
				<el-input v-model="query.name" :placeholder="$t('job.name_form')" clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('job.run_on')">
				<el-input v-model="query.node_uuid" :placeholder="$t('job.uuid_form')" clearable></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-search" @click="getTableData">{{$t('common.search')}}</el-button>
				<el-button icon="el-icon-search" @click="handleReset">{{$t('common.reset')}}</el-button>
			</el-form-item>
		</el-form>

		<div class="log-manager__table">
			<el-table
				v-loading="tableLoading"
				:data="logList"
				:row-class-name="tableRowClassName"
				@selection-change="handleSelectedRows"
				style="width: 100%"

			>
			  <el-table-column type="expand">
				<template slot-scope="scope">
					<el-form label-position="left" inline class="job-table-expand">
						<el-form-item :label="$t('log.id')">
							<span>{{ scope.row.id }}</span>
						</el-form-item>
						<el-row :gutter="20">
							<el-col :xs="24" :md="12">
								<el-form-item :label="$t('job.id')">
									<span>{{ scope.row.job_id }}</span>
								</el-form-item>
							</el-col>
							<el-col :xs="24" :md="12">	
								<el-form-item :label="$t('job.name')">
									<span>{{ scope.row.name }}</span>
								</el-form-item>
							</el-col>
						</el-row>
						<el-form-item :label="$t('job.run_on')">
									<span>{{ scope.row.node_uuid }}</span>
							</el-form-item>
					
						<el-form-item :label="$t('log.execution_time')">
									<span>{{formatDate( scope.row.start_time )}} ---> {{scope.row.end_time?formatDate(scope.row.end_time):0}}</span>
								</el-form-item>
						<el-row :gutter="20">
							<el-col :xs="24" :md="12">
								<el-form-item :label="$t('log.status')" label-width="120px">
											<el-tag  effect="plain" :type="scope.row.success?`success`:(scope.row.end_time?'danger':'info')">
													{{ scope.row.success?`Success`:(scope.row.end_time?'Failed':"Running") }}
											</el-tag>
								</el-form-item>


							</el-col>
							<el-col :xs="24" :md="12">
								<el-form-item :label="$t('job.retry_time')">
									<span>{{ scope.row.retry_times }}</span>
								</el-form-item>
							</el-col>
							
						</el-row>
					
					<el-form-item :label="$t('log.output')">
						<textarea  class="expand-text" v-model="scope.row.output"></textarea> 
					</el-form-item>
					</el-form>
				</template>
			  </el-table-column>
				<el-table-column prop="job_id"  :label="$t('job.id')"></el-table-column>
				<el-table-column prop="name" :label="$t('job.name')" ></el-table-column>
				<el-table-column prop="success" :label="$t('log.status')">
					<template slot-scope="scope">
						<el-tag  effect="plain" :type="scope.row.success?`success`:(scope.row.end_time?'danger':'info')">
								{{ scope.row.success?`Success`:(scope.row.end_time?'Failed':"Running") }}
						</el-tag>
						<el-button v-if="scope.row.success==0&&scope.row.end_time>0" @click="onceExecute(scope.row)" style="margin-left:5px; color:blue;font-size:15px;" type="text">
							<i  class="el-icon-refresh-left"></i>
						</el-button>
						<!-- <i v-if="scope.row.success==0&&scope.row.end_time==0" @click="kill(scope.row)" style="margin-left:5px; color:blue;font-size:15px;" class="el-icon-video-pause"></i> -->
					</template>
				</el-table-column>
				<el-table-column prop="retry_times" :label="$t('job.retry_time')"></el-table-column>
				<el-table-column prop="start_time" :label="$t('log.execution_time')" width="320px">
						<template slot-scope="scope">
							<span>{{formatDate( scope.row.start_time )}}/ {{scope.row.end_time?formatDate(scope.row.end_time):0}}</span>
						</template>
				</el-table-column>
			</el-table>
		</div>

		<pagination
			:total="total"
			:page-number.sync="query.page"
			:page-size.sync="query.page_size"
			@pagination="getTableData"
		></pagination>

	</div>
</template>

<script>
/**
 * 日志管理
 */
import _ from 'lodash'
import { scroll } from '@/utils/core'
import tableMng from '@/utils/tableMng'
import { exportExcel } from '@/utils/excle'
import { formatDate } from '@/utils/date'

const defaultQuery = {
	name: '',
	job_id: null,
	node_uuid: '',
	success:null,
	page: 1,
	page_size: 10,
}

export default {
	name: 'Log',
	props: ['jid'],
	data() {
		return {
			tableMng,
			logList: [],
			query: _.cloneDeep(defaultQuery),
			total: 0,
			selectedRows: [],
			tableLoading: false,
			exportLoading: false,
		}
	},
	created() {
		this.getTableData()
	},
	methods: {
		formatDate(time) {
			time = time * 1000
			let date = new Date(time)
			return formatDate(date, 'yyyy-MM-dd hh:mm:ss')
      },
		tableRowClassName({row, rowIndex}) {
			if (row.success === false) {
			return 'warning-row';
			} else {
			return 'success-row';
			}
		},
		//获取日志列表
		async getTableData() {
			this.tableLoading = true
			if(this.jid){
				this.query.job_id=Number(this.jid);
			}
      if(this.query.job_id==''){
        this.query.job_id=null
      }
			const data = await this.$api.job.getLogList(this.query)
			this.logList = data.list
			this.total = data.total
			this.tableLoading = false
			const scrollElement = document.querySelector('.inner-layout__page')
			scroll(scrollElement, 0, 800)
		},
		// 重置查询
		handleReset() {
			this.query = _.cloneDeep(defaultQuery)
			this.jid=null
			this.getTableData()
		},
		// 多选
		handleSelectedRows(rows) {
			this.selectedRows = rows
		},
		handleFilterGender(value, row, column) {
			const property = column['property']
			return row[property] === value
		},
		//立即执行
		async  onceExecute(row){
			const res=await this.$api.job.once({ job_id:row.job_id,node_uuid:row.node_uuid })
			console.log(res);
			this.getTableData()			
		},
		//停止
		async  kill(row){
			await this.$api.job.kill({ job_id:row.job_id ,node_uuid:row.node_uuid})
		}
	},
}
</script>

<style lang="scss" >
.log-manager {
	background-color: #fff;
	padding: 1em;

	.log-manager__header {
		display: flex;
		flex-wrap: wrap;
		align-items: center;
		justify-content: space-between;

		.title,
		.operation {
			margin-bottom: 1em;
		}
	}
}
.el-table .warning-row {
    background: oldlace;
}

.el-table .success-row {
    background: #f0f9eb;
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
.expand-text{
	font-size: 10px;
	width: 450px;
	color: #fff;
	background-color: rgb(14, 13, 13);
	resize:none;
	border-radius: 5px;
}
</style>
