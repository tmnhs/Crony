<template>
	<div class="script-list" >
		<div class="script-list__header">
			<div class="title">
				<section-title :name="$t('script.script_list')" />
			</div>

			<div class="operation">
				<el-button type="primary" icon="el-icon-plus" @click="handleAdd">{{$t('common.add')}}</el-button>
				<el-button type="danger" icon="el-icon-minus" @click="handleDelete(null)" >{{$t('common.deletes')}}</el-button>
				<el-button icon="el-icon-download" :loading="exportLoading" @click="handleExport">{{$t('common.export')}}</el-button>

			</div>
		</div>

		<el-form :inline="true" :model="query">
			<el-form-item :label="$t('script.id')">
				<el-input v-model.number="query.id" :placeholder="$t('script.id_form')" clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('script.name')">
				<el-input v-model="query.name" :placeholder="$t('script.name_form')" clearable></el-input>
			</el-form-item>

			<el-form-item>
				<el-button type="primary" icon="el-icon-search" @click="getScriptList"> {{$t('common.search')}} </el-button>
				<el-button icon="el-icon-search" @click="handleReset">{{$t('common.reset')}}</el-button>

			</el-form-item>
		</el-form>

		<div class="script-list__table">
			<el-table
				:data="scriptList"
				v-loading="scriptTableLoading"
				@selection-change="handleSelectedRows"
			>	
				<el-table-column type="selection" width="50"></el-table-column>

				<el-table-column prop="id" :label="$t('script.id')" width="120px"></el-table-column>
				<el-table-column prop="name" :label="$t('script.name')" width="180px" :show-overflow-tooltip="true">
        </el-table-column>
				<el-table-column  :label="$t('script.command')" min-width="500px">
          <template  slot-scope="scope">
            <textarea  class="command-text" v-model="scope.row.command"></textarea>
          </template>
        </el-table-column>
        <el-table-column prop="created" width="240px" :label="$t('script.create_time')" sortable>
          <template slot-scope="scope">
            {{ formatDate(scope.row.created)}}
          </template>
        </el-table-column>
				<el-table-column label="操作" width="200px">
					<template slot-scope="scope">

						<router-link :to="`/script/edit/${scope.row.id}`">
							<el-button type="text">{{$t('common.edit')}}</el-button>
						</router-link>
					
						<el-divider direction="vertical"></el-divider>
						<el-button type="text" @click="handleDelete(scope.row)"> {{$t('common.delete')}} </el-button>
					</template>
				</el-table-column>
			</el-table>
		</div>

		<pagination
			:total="total"
			:page-number.sync="query.page"
			:page-size.sync="query.page_size"
			@pagination="getScriptList"
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
	page: 1,
	page_size: 10,
}
export default {
	name: 'ScriptList',
	data() {
		return {
			tableMng,
			userLoading: false,
			scriptList: [],
			scriptTableLoading: false,
			query: _.cloneDeep(defaultQuery),
			total: 0,
			selectedRows: [],
			exportLoading: false,
		}
	},
	created() {
		this.getScriptList()
	},
	methods: {
		formatDate(time) {
			time = time * 1000
			let date = new Date(time)
			return formatDate(date, 'yyyy-MM-dd hh:mm:ss')
      },
		// 获取任务列表
		async getScriptList() {
			this.scriptTableLoading = true
      if(this.query.id==''){
        this.query.id=null
      }
			const data = await this.$api.script.getScriptList(this.query)
			this.scriptList = data.list
			this.total = data.total
			this.scriptTableLoading = false
			const scrollElement = document.querySelector('.inner-layout__page')
			scroll(scrollElement, 0, 800)
		},
		// 跳转到新建任务页面
		handleAdd() {
			this.$router.push('/script/add')
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
						await this.$api.script.removeScript({ ids:ids })
						this.$message.success(i18n.t('common.delete_success'))
						this.getScriptList()
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
			const header = ['ID', 'Name','Command', 'CratedTime']
			const res = await this.$api.script.getScriptList()
			const data = res.list.map((item, index) => {
				return {
						id: item.id,
						name: item.name,
						command: item.command,
						created:item.created,
				}
			})
			exportExcel(header, data, 'script Data Sheet')
			this.exportLoading = false
		},
		
		// 重置查询
		handleReset() {
			this.query = _.cloneDeep(defaultQuery)
			this.getScriptList()
		},
	},
}
</script>

<style lang="scss" >
.script-list {
	min-height: 100%;
	padding: 1em;
	box-sizing: border-box;
	background-color: #fff;

	.script-list__header {
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
.script-table-expand label{
	font-size: 12px;
	width: 90px;
    color: #99a9bf;
}
.script-table-expand {
	
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
	width: 500px;
  height: 100px;
	color: #fff;
	background-color: rgb(14, 13, 13);
	resize:none;
	border-radius: 5px;
}

</style>
