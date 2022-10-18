<template>
	<div>
        <el-card>
            <div class="card-top">
                <div class="card-top-title">{{$t('home.hello')}}</div>
                <div class="card-top-dot"> {{$t('home.weather')}}</div>
            </div>

            <img src="http://qny.tmnhs.top/crony/guanliyuan.png" alt="" class="card-top-img">
        </el-card>
        <el-card class="box-card" style="margin:20px 0;">
		    <statistics style="margin:50px 20px" :statistics="statis" @onChangeType="hanldeChangeType"></statistics>
        </el-card>
		<el-row >
			<el-col :lg="11" :sm="24" >
                <div class="card-bottom">
				    <p class="card-bottom-title">{{$t('home.today_job')}}</p>
                </div>
				<div :id="chartId"></div>
			</el-col>
			<el-col  :lg="11" :sm="24" style="float:right;">
                 <div class="card-bottom">
				    <p class="card-bottom-title">{{$t('home.week_job')}}</p>
                </div>
				<line-chart :statistic-type="statisticType"></line-chart>
			</el-col>
	
		</el-row>
		
	</div>
</template>

<script>
import Statistics from './components/statistics'
import LineChart from './components/LineChart'
import G2 from "@antv/g2";
import {
	DataSet
} from '@antv/data-set';
const defaultStatis={
		normal_node_count:0,
		fail_node_count:0,
		job_exc_success_count:0,
		job_running_count:0,
		job_exc_fail_count:0,
}

export default {
	name: 'Dashboard',
	components: {
		Statistics,
		LineChart,
	},
	data() {
		return {
			statisticType: 'visite',
			statis:{...defaultStatis},
			chartData: [],
			chartId: "chart" + +new Date() + ((Math.random() * 1000).toFixed(0) + ""),
			height:360,
		}
	},
	mounted() {
         this.createChart(this.chartId);
	},
	methods: {
		hanldeChangeType(type) {
			this.statisticType = type
		},
		async createChart(container) {
                const res = await this.$api.dashboard.getTodayStatics()
				this.statis={
						fail_node_count: res.fail_node_count,
						job_exc_fail_count:res.job_exc_fail_count,
						job_exc_success_count:res.job_exc_success_count,
						job_running_count:res.job_running_count,
						normal_node_count:res.normal_node_count,
					}
                const data=[
                    {
                    type: '执行成功数',
                    value: this.statis.job_exc_success_count
                  }, {
                    type: '执行失败数',
                    value:this.statis.job_exc_fail_count
                  }, {
                    type: '运行中',
                    value: this.statis.job_running_count
                  }
                  ]
                const total=this.statis.job_exc_success_count+this.statis.job_exc_fail_count+this.statis.job_running_count
                let startAngle = -Math.PI / 2 - Math.PI / 4;
                let ds = new DataSet();
                let dv = ds.createView().source(data);
                dv.transform({
                    type: 'percent',
                    field: 'value',
                    dimension: 'type',
                    as: 'percent'
                });

                let chart = new G2.Chart({
                    container: container,
                    forceFit: true,
                    height: this.height,
                    padding: [50, 0],
                    background: {
                        fill: "#fff"
                    }
                })
                chart.source(dv);
                chart.legend(false);
                chart.tooltip({
                    showTitle: false
                })
                chart.coord('theta', {
                    radius: 0.75,
                    innerRadius: 0.5,
                    startAngle: startAngle,
                    endAngle: startAngle + Math.PI * 2
                });
                chart.intervalStack().position('value').color('type', ['#06d6a0', '#ef476f', '#ffd166', '#2295ff',
                    '#48adff',
                    '#6fc3ff', '#96d7ff', '#bde8ff'
                ]).opacity(1).label('percent', {
                    offset: -20,
                    textStyle: {
                        fill: 'white',
                        fontSize: 12
                    },
                    formatter: function formatter(val) {
                        return parseInt(val * 100) + '%';
                    }
                });
                chart.guide().html({
                    position: ['50%', '50%'],
                    html: '<div class="g2-guide-html"><p class="title">总计'+ total+'</div>'
                });
                chart.render();

                //绘制label
                let OFFSET = 20;
                let APPEND_OFFSET = 100;
                let LINEHEIGHT = 60;
                let coord = chart.get('coord'); // 获取坐标系对象
                let center = coord.center; // 极坐标圆心坐标
                let r = coord.radius; // 极坐标半径
                let canvas = chart.get('canvas');
                let canvasWidth = chart.get('width');
                let canvasHeight = chart.get('height');
                let labelGroup = canvas.addGroup();
                let labels = [];
                addPieLabel(chart);
                canvas.draw();
                chart.on('afterpaint', function () {
                    addPieLabel(chart);
                });
                //main
                function addPieLabel() {
                    let halves = [
                        [],
                        []
                    ];
                    let data = dv.rows;
                    let angle = startAngle;
                    for (let i = 0; i < data.length; i++) {
                        let percent = data[i].percent;
                        let targetAngle = angle + Math.PI * 2 * percent;
                        let middleAngle = angle + (targetAngle - angle) / 2;
                        angle = targetAngle;
                        let edgePoint = getEndPoint(center, middleAngle, r);
                        let routerPoint = getEndPoint(center, middleAngle, r + OFFSET);
                        //label
                        let label = {
                            _anchor: edgePoint,
                            _router: routerPoint,
                            _data: data[i],
                            x: routerPoint.x,
                            y: routerPoint.y,
                            r: r + OFFSET,
                            fill: '#bfbfbf'
                        };
                        // 判断文本的方向
                        if (edgePoint.x < center.x) {
                            label._side = 'left';
                            halves[0].push(label);
                        } else {
                            label._side = 'right';
                            halves[1].push(label);
                        }
                    } // end of for
                    let maxCountForOneSide = parseInt(canvasHeight / LINEHEIGHT, 10);
                    halves.forEach(function (half, index) {
                        // step 2: reduce labels
                        if (half.length > maxCountForOneSide) {
                            half.sort(function (a, b) {
                                return b._percent - a._percent;
                            });
                            half.splice(maxCountForOneSide, half.length - maxCountForOneSide);
                        }
                        // step 3: distribute position (x and y)
                        half.sort(function (a, b) {
                            return a.y - b.y;
                        });
                        antiCollision(half, index);
                    });
                }

                function getEndPoint(center, angle, r) {
                    return {
                        x: center.x + r * Math.cos(angle),
                        y: center.y + r * Math.sin(angle)
                    };
                }

                function drawLabel(label) {
                    let _anchor = label._anchor,
                        _router = label._router,
                        fill = label.fill,
                        y = label.y;
                    let labelAttrs = {
                        y: y,
                        fontSize: 12, // 字体大小
                        fill: '#808080',
                        text: label._data.type + '\n' + label._data.value,
                        textBaseline: 'bottom'
                    };
                    let lastPoint = {
                        y: y
                    };
                    if (label._side === 'left') {
                        // 具体文本的位置
                        lastPoint.x = APPEND_OFFSET;
                        labelAttrs.x = APPEND_OFFSET; // 左侧文本左对齐并贴着画布最左侧边缘
                        labelAttrs.textAlign = 'left';
                    } else {
                        lastPoint.x = canvasWidth - APPEND_OFFSET;
                        labelAttrs.x = canvasWidth - APPEND_OFFSET; // 右侧文本右对齐并贴着画布最右侧边缘
                        labelAttrs.textAlign = 'right';
                    }
                    // 绘制文本
                    let text = labelGroup.addShape('Text', {
                        attrs: labelAttrs
                    });
                    labels.push(text);
                    // 绘制连接线
                    let points = void 0;
                    if (_router.y !== y) {
                        // 文本位置做过调整
                        points = [
                            [_anchor.x, _anchor.y],
                            [_router.x, y],
                            [lastPoint.x, lastPoint.y]
                        ];
                    } else {
                        points = [
                            [_anchor.x, _anchor.y],
                            [_router.x, _router.y],
                            [lastPoint.x, lastPoint.y]
                        ];
                    }
                    labelGroup.addShape('polyline', {
                        attrs: {
                            points: points,
                            lineWidth: 1,
                            stroke: fill
                        }
                    });
                }

                function antiCollision(half, isRight) {
                    let startY = center.y - r - OFFSET - LINEHEIGHT;
                    let overlapping = true;
                    let totalH = canvasHeight;
                    let i = void 0;
                    let maxY = 0;
                    let minY = Number.MIN_VALUE;
                    let boxes = half.map(function (label) {
                        let labelY = label.y;
                        if (labelY > maxY) {
                            maxY = labelY;
                        }
                        if (labelY < minY) {
                            minY = labelY;
                        }
                        return {
                            size: LINEHEIGHT,
                            targets: [labelY - startY]
                        };
                    });
                    if (maxY - startY > totalH) {
                        totalH = maxY - startY;
                    }
                    while (overlapping) {
                        boxes.forEach(function (box) {
                            let target = (Math.min.apply(minY, box.targets) + Math.max.apply(minY, box.targets)) /
                                2;
                            box.pos = Math.min(Math.max(minY, target - box.size / 2), totalH - box.size);
                        });
                        // detect overlapping and join boxes
                        overlapping = false;
                        i = boxes.length;
                        while (i--) {
                            if (i > 0) {
                                let previousBox = boxes[i - 1];
                                let box = boxes[i];
                                if (previousBox.pos + previousBox.size > box.pos) {
                                    // overlapping
                                    previousBox.size += box.size;
                                    previousBox.targets = previousBox.targets.concat(box.targets);
                                    // overflow, shift up
                                    if (previousBox.pos + previousBox.size > totalH) {
                                        previousBox.pos = totalH - previousBox.size;
                                    }
                                    boxes.splice(i, 1); // removing box
                                    overlapping = true;
                                }
                            }
                        }
                    }
                    // step 4: normalize y and adjust x
                    i = 0;
                    boxes.forEach(function (b) {
                        let posInCompositeBox = startY; // middle of the label
                        b.targets.forEach(function () {
                            half[i].y = b.pos + posInCompositeBox + LINEHEIGHT / 2;
                            posInCompositeBox += LINEHEIGHT;
                            i++;
                        });
                    });
                    // (x - cx)^2 + (y - cy)^2 = totalR^2
                    half.forEach(function (label) {
                        let rPow2 = label.r * label.r;
                        let dyPow2 = Math.pow(Math.abs(label.y - center.y), 2);
                        if (rPow2 < dyPow2) {
                            label.x = center.x;
                        } else {
                            let dx = Math.sqrt(rPow2 - dyPow2);
                            if (!isRight) {
                                // left
                                label.x = center.x - dx;
                            } else {
                                // right
                                label.x = center.x + dx;
                            }
                        }
                        drawLabel(label);
                    });
                }

            },
		
	},
}


</script>


<style scoped>
.card-bottom{
	background: #fff;
    height: 64px;
}
.card-bottom-title{
	font-weight: 800;
	font-size: 20px;
	width: 100%;
	text-align: center;
	margin: 0 auto;
    line-height: 64px;
}
.card-top-title{
    font-size: 22px;
    color: #343844;
    font-weight: bolder;
}
.card-top{
    float: left;

}
.card-top-img{
    font-size: 22px;
    color: #343844;
    height: 200px;
    margin-right: 50px;
    width:200px;
    float: right;

}
.card-top-dot{
     font-size: 16px;
    color: #6B7687;
    margin-top: 24px;
}
</style>