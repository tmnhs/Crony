<template>
	<div class="user-manager">
		<div class="user-manager__header">
			<div class="title">
				<section-title :name="$t('user.user_list')" />
			</div>
			<div class="operation">
				<el-button type="primary" icon="el-icon-plus" @click="handleEdit({})"> {{$t('common.add')}} </el-button>
				<el-button type="danger" icon="el-icon-minus" @click="handleDelete({})" > {{$t('common.deletes')}} </el-button>
			</div>
		</div>

		<el-form ref="queryForm" :inline="true" :model="query">
			<el-form-item label="ID">
				<el-input v-model.number="query.id"  :placeholder="$t('common.id_query')" clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('user.name')">
				<el-input v-model="query.username" :placeholder="$t('user.name_query')" clearable></el-input>
			</el-form-item>
			<el-form-item :label="$t('user.role')">
				<el-select v-model.number="query.role" :placeholder="$t('user.role_query')"  clearable>
					<el-option
						v-for="item in tableMng.getTable('role')"
						:key="item.id"
						:label="item.name"
						:value="item.id"
					></el-option>
				</el-select>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-search" @click="getTableData">{{$t('common.search')}}</el-button>
				<el-button icon="el-icon-search" @click="handleReset">{{$t('common.reset')}}</el-button>
			</el-form-item>
		</el-form>

		<div class="user-manager__table">
			<el-table
				v-loading="tableLoading"
				:data="userList"
			
				@selection-change="handleSelectedRows"
			>
				<el-table-column type="selection" width="50px"></el-table-column>
				<el-table-column prop="id" label="ID"></el-table-column>
				<el-table-column prop="username" :label="$t('user.name')"></el-table-column>
				<el-table-column prop="email" :label="$t('user.email')"></el-table-column>
				
				<el-table-column prop="role" :label="$t('user.role')">
					<template slot-scope="scope">
						{{ tableMng.getNameById('role', scope.row.role) }}
					</template>
				</el-table-column>
				<el-table-column prop="created" :label="$t('user.created')" sortable>
					<template slot-scope="scope">
						{{ formatDate(scope.row.created)}}
					</template>
				</el-table-column>
				<el-table-column :label="$t('common.operation')" width="120px">
					<template slot-scope="scope">
						<el-button type="text" @click="handleEdit(scope.row)" >{{$t('common.edit')}}</el-button>
						<el-divider direction="vertical"></el-divider>
						<el-button type="text" @click="handleDelete(scope.row)" >{{$t('common.delete')}}</el-button>
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
		<user-edit ref="userEdit" @success="getTableData"></user-edit>
	</div>
</template>

<script>
/**
 * 用户管理
 */
import _ from 'lodash'
import { scroll } from '@/utils/core'
import tableMng from '@/utils/tableMng'
import { formatDate } from '@/utils/date'
import UserEdit from './components/UserEdit'
import i18n from  "@/assets/lang";
const defaultQuery = {
	id:null,
	username: '',
	email: '',
	role: null ,
	page: 1,
	page_size: 10,
}

export default {
	name: 'User',
	components: {
		UserEdit,
	},
	data() {
		return {
			tableMng,
			userList: [],
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
		//获取用户列表
		async getTableData() {
			// this.tableLoading = true
      if(this.query.id==''){
        this.query.id=null
      }
			const data = await this.$api.user.getList(this.query)
			this.userList = data.list
			this.total = data.total
			// this.tableLoading = false
			const scrollElement = document.querySelector('.inner-layout__page')
			scroll(scrollElement, 0, 800)
		},
		// 重置查询
		handleReset() {
			this.query = _.cloneDeep(defaultQuery)
			this.getTableData()
		},
		// 编辑/新增
		handleEdit(row) {
			this.$refs.userEdit.open(row)
		},
		// 删除
		handleDelete(row) {
			let ids = []
			let usernames = []
			if (row.id) {
				ids = [row.id]
				usernames = [row.username]
			} else {
				ids = this.selectedRows.map(row => row.id)
				usernames = this.selectedRows.map(row => row.username)
			}
			if (usernames.length === 0) {
				this.$message.warning(i18n.t('common.delete_warn'))
			} else {
				this.$confirm(i18n.t('common.delete_question')+`：“${usernames.join('，')}”？`, i18n.t('common.prompt'), {
					type: 'warning',
				})
					.then(async () => {
						const res=await this.$api.user.remove({ids:ids})
						this.$message.success(i18n.t('common.delete_success'))
						this.getTableData()
					})
					.catch(() => {})
			}
		},
		// 多选
		handleSelectedRows(rows) {
			this.selectedRows = rows
		},
		handleFilterGender(value, row, column) {
			const property = column['property']
			return row[property] === value
		},
		
		formatDate(time) {
			time = time * 1000
			let date = new Date(time)
			return formatDate(date, 'yyyy-MM-dd hh:mm:ss')
      }
	},
}
</script>

<style lang="scss" scoped>
.user-manager {
	background-color: #fff;
	padding: 1em;

	.user-manager__header {
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
</style>
