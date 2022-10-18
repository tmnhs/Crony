<template>
	<div class="node-manager">
		<el-row>
				<div class="node-manager__header">
			<div class="title">
				<section-title :name="$t('node.node_list')" />
			</div>
		</div>
		</el-row>
		<el-form ref="queryForm" :inline="true" :model="query">
			<el-form-item label="UUID:">
				<el-input v-model="query.uuid"  :placeholder="$t('node.uuid_form')"  clearable></el-input>
			</el-form-item>
			<el-form-item label="IP:">
				<el-input v-model="query.ip" :placeholder="$t('node.ip_form')" clearable></el-input>
			</el-form-item>
			<el-form-item>
				<el-button type="primary" icon="el-icon-search" @click="getNodeList">{{$t('common.search')}}</el-button>
				<el-button icon="el-icon-search" @click="handleReset">{{$t('common.reset')}}</el-button>
			</el-form-item>
		</el-form>
		<el-row :gutter="30">
			<el-col :lg="6" :sm="12" v-for="node in nodeList" :key="node.id">
				<div class="dashboard-statistic " :class="node.status==1?`node--blue`:`node--red`">
					<div class="node-top">
						<p class="node-name">{{ node.ip }}:{{node.pid}}</p>
						<div class="node-uuid">
							<p  ><span style="font-weight:400;">UUID:</span>{{ node.uuid }}  <span  @click="handleCopyText(node.uuid)"><i   class="el-icon-document-copy"></i></span></p>
							<p ><span style="font-weight:400;">HostName: </span>{{ node.hostname }}</p>
							<p ><span style="font-weight:400;">{{$t('node.created')}}  :</span>{{ formatDate(node.up) }} <span style="font-weight:400;"> {{$t('node.version')}}: </span>{{ node.version }}</p>
							<div class="node-butttom">
									<button v-if="node.status==2" type="successs" @click="handleDeleteNode(node.uuid)" class="node-button small"> {{$t('common.delete')}}</button>
									<button type="successs" @click="handleNodeJob(node.uuid)" class="node-button small"> {{$t('node.job_count')}}:{{ node.job_count }}</button>
									<button type="successs" @click="handleSystemInfo(node.uuid)" class="node-button small" > {{$t('node.server')}} </button>
							</div>
		
	
						</div>
					</div>
				</div>

			</el-col>
	</el-row>

	</div>
</template>

<script>
/**
 * 节点管理
 */
import _ from 'lodash'
import { formatDate } from '@/utils/date'
const defaultQuery = {
	ip: '',
	uuid: '',
	up: null,
	page: 1,
	status:null,
	page_size: 20,
}

export default {
	name: 'Node',
	data() {
		return {
			query: _.cloneDeep(defaultQuery),
			total: 0,
			selectedRows: [],
			exportLoading: false,
			nodeList: [],
		}
	},
	created() {
		this.getNodeList()
	},
	methods: {
		formatDate(time) {
			time = time * 1000
			let date = new Date(time)
			return formatDate(date, 'yyyy-MM-dd hh:mm:ss')
      },
	  // 重置查询
		handleReset() {
			this.query = _.cloneDeep(defaultQuery)
			this.getNodeList()
		},
		//获取节点列表
		async getNodeList() {
			const data = await this.$api.node.getNodeList(this.query)
			this.nodeList = data.list
			this.total = data.total
		},
		async handleDeleteNode(uuid){
			const data = await this.$api.node.remove({uuid:uuid})
			this.getNodeList()
		},
		handleCopyText(text) {
			const wrap = document.createElement('p')
			wrap.innerText = text
			document.body.appendChild(wrap)
			this.handleCopyElem(wrap)
			document.body.removeChild(wrap)
		},
		handleCopyNode() {
			this.handleCopyElem(this.$refs.content)
		},
		handleSystemInfo(uuid) {
			this.$router.push('/node/system/'+uuid)
		},
		handleNodeJob(uuid){
			this.$router.push('/node/job/'+uuid)
		},
		handleCopyElem(elem) {
			if (!elem) {
				return this.$message.warning('there is nothing to copy')
			}
			const selection = window.getSelection()
			// 如果剪切板中已经有复制了的内容，需要清掉。
			if (selection.rangeCount > 0) {
				selection.removeAllRanges()
			}
			const range = document.createRange()
			range.selectNodeContents(elem)
			selection.addRange(range)
			const result = document.execCommand('Copy')
			// 清除选中的内容,也可以使用 window.getSelection().removeAllRanges()
			range.collapse(false)
			if (result) {
				this.$message.success('copy success')
			} else {
				this.$message.error('copy fail')
			}
		},
		


	},
}
</script>

<style lang="scss" scoped>
.node-manager {
	background-color: #fff;
	padding: 1em;

	.node-manager__header {
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
.el-card {
	min-width: 200px;
	margin-right: 10px;
	border-radius: 20px;
}
.dashboard-statistic {
	display: flex;
	position: relative;
	height: 120px;
	// color: #fff;
	box-shadow: 0px 0px 10px rgba(100, 100, 100, 0.5);
	cursor: pointer;
	margin-bottom: 10px;
	border-radius: 10px;

	.node-top {
		padding: 2px;
		width: 100%;
		p {
			color: #fff !important;

		}

		.node-name {
			font-size: 20px;
			border-bottom: 2px solid #24d3a4;
			font-weight: bolder;
		}
		.node-uuid{
			margin-top: 5px;
			font-size: 5px;
			font-weight: 100;
		} 
	}
	.node-butttom{
		margin-top:2px ;
		 float: right;
	}
}

.node {
	&--green {
		background-color: #06d6a0;
	}

	&--yellow {
		background-color: #ffd166;
	}

	&--blue {
		background-color: #06aed5;
	}

	&--red {
		background-color: #ef476f;
	}
}

.node-button{
	margin-top: 5px;
	margin-left: 100px;
	// font-weight: bolder;
	// float: right;

	display: inline-block;
	*display: inline;
	zoom: 1;
	padding: 6px 20px;
	margin: 0;
	cursor: pointer;
	border: 1px solid #bbb;
	overflow: visible;
	font: bold 13px arial, helvetica, sans-serif;
	text-decoration: none;
	white-space: nowrap;
	color: #06aed5;
	
	background-color: #ddd;
	background-image: -webkit-gradient(linear, left top, left bottom, from(rgba(255,255,255,1)), to(rgba(255,255,255,0)));
	background-image: -webkit-linear-gradient(top, rgba(255,255,255,1), rgba(255,255,255,0));
	background-image: -moz-linear-gradient(top, rgba(255,255,255,1), rgba(255,255,255,0));
	background-image: -ms-linear-gradient(top, rgba(255,255,255,1), rgba(255,255,255,0));
	background-image: -o-linear-gradient(top, rgba(255,255,255,1), rgba(255,255,255,0));
	background-image: linear-gradient(top, rgba(255,255,255,1), rgba(255,255,255,0));
	
	-webkit-transition: background-color .2s ease-out;
	-moz-transition: background-color .2s ease-out;
	-ms-transition: background-color .2s ease-out;
	-o-transition: background-color .2s ease-out;
	transition: background-color .2s ease-out;
	background-clip: padding-box; /* Fix bleeding */
	-moz-border-radius: 3px;
	-webkit-border-radius: 3px;
	border-radius: 3px;
	-moz-box-shadow: 0 1px 0 rgba(0, 0, 0, .3), 0 2px 2px -1px rgba(0, 0, 0, .5), 0 1px 0 rgba(255, 255, 255, .3) inset;
	-webkit-box-shadow: 0 1px 0 rgba(0, 0, 0, .3), 0 2px 2px -1px rgba(0, 0, 0, .5), 0 1px 0 rgba(255, 255, 255, .3) inset;
	box-shadow: 0 1px 0 rgba(0, 0, 0, .3), 0 2px 2px -1px rgba(0, 0, 0, .5), 0 1px 0 rgba(255, 255, 255, .3) inset;
	text-shadow: 0 1px 0 rgba(255,255,255, .9);
	
	-webkit-touch-callout: none;
	-webkit-user-select: none;
	-khtml-user-select: none;
	-moz-user-select: none;
	-ms-user-select: none;
	user-select: none;
}
.small{
	 padding: 4px 12px;
}
</style>
