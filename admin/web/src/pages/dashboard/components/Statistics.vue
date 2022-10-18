<template>
	<el-row :gutter="30">
		<el-col :lg="8" :sm="12">
		<router-link to="/node/index">
				<div class="dashboard-statistic grid--blue" >
				<div class="grid-left">
					<base-icon name="view" icon-class="grid-icon" :size="48"></base-icon>
				</div>

				<div class="grid-right">
					<p class="grid-name">{{$t('home.node')}}</p>
					<countTo class="grid-number" :startVal="0" :endVal="statistics.normal_node_count" :duration="3000"></countTo>
					<span class="grid-number">/</span>
					<countTo class="grid-number" :startVal="0" :endVal="statistics.fail_node_count" :duration="3000"></countTo>
				</div>
			</div>
		</router-link>
		</el-col>
		<el-col :lg="8" :sm="12">
			<router-link to="/job/list/">
				<div class="dashboard-statistic grid--green" >
				<div class="grid-left">
					<base-icon name="view" icon-class="grid-icon" :size="48"></base-icon>
				</div>

				<div class="grid-right">
					<p class="grid-name">{{$t('home.job')}}</p>
					<countTo class="grid-number" :startVal="0" :endVal="statistics.job_exc_success_count+statistics.job_running_count+statistics.job_exc_fail_count" :duration="3000"></countTo>
					<span class="grid-number">/</span>
					<countTo class="grid-number" :startVal="0" :endVal="statistics.job_running_count" :duration="3000"></countTo>
				</div>
			</div>
			</router-link>
		</el-col>
		<el-col :lg="8" :sm="12">
			<router-link to="/log/index">
					<div class="dashboard-statistic grid--red" >
				<div class="grid-left">
					<base-icon name="view" icon-class="grid-icon" :size="48"></base-icon>
				</div>

				<div class="grid-right">
					<p class="grid-name">{{$t('home.job_fail')}}</p>
					<countTo class="grid-number" :startVal="0" :endVal="statistics.job_exc_fail_count" :duration="3000"></countTo>
				</div>
			</div>
			</router-link>
		</el-col>
	</el-row>
</template>

<script>
import countTo from 'vue-count-to'

export default {
	props: ["statistics"],
	components: { countTo },
	data() {
		return {
			
		}
	},
}
</script>

<style lang="scss" scoped>
.dashboard-statistic {
	display: flex;
	position: relative;
	height: 120px;
	color: #fff;
	border-radius: 5px;
	box-shadow: 0px 0px 10px rgba(100, 100, 100, 0.5);
	cursor: pointer;
	margin-bottom: 20px;

	&::before {
		content: '';
		position: absolute;
		top: 0px;
		height: 100%;
		width: 100px;
		background-color: rgba(0, 0, 0, 0.05);
		-webkit-transition: all 0.95s;
		transition: width 0.5s;
	}

	&:hover::before {
		width: 100%;
		box-sizing: border-box;
	}

	.grid-left {
		width: 100px;
		line-height: 108px;
		vertical-align: top;
		text-align: center;

		.grid-icon {
			vertical-align: middle;
		}
	}

	.grid-right {
		padding: 16px;

		p {
			color: #fff !important;
		}

		.grid-name {
			font-size: 20px;
			margin-bottom: 10px;
		}

		.grid-number {
			font-size: 30px;
		}
	}
}

.grid {
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
</style>
