<template>
  <div :id="chartId"></div>
</template>

<script>
  import G2 from '@antv/g2';

  export default {
    props: {
      height: {
        type: Number,
        default: 360
      }
    },
    data() {
      return {
        chartId: 'chart' + +new Date() + ((Math.random() * 1000).toFixed(0) + ''),
	  }
    },

    mounted() {
      this.createChart(this.chartId);
    },
    methods: {

      async createChart(container) {
		const res = await this.$api.dashboard.getWeekStatics()
		var data = []
		res.success_date_count.forEach(item => {
			data.push({
				date:item.date,
				count:item.count,
				type:"success",
				color:"green",
			})
		});

		res.fail_date_count.forEach(item => {
		data.push({
			date:item.date,
			count:item.count,
			type:"fail",
			color:"red",
		})
		});
		data.sort(function(a, b) {
			return a.count-b.count
		})
        let chart = new G2.Chart({
          container: container,
          forceFit: true,
          height: this.height,
          background: {
            fill: "#fff"
          }
        });
        chart.source(data);
        chart.scale({
          date: {
            range: [0, 1],
          }
        })
        chart.axis('count', {
          label: { formatter: val => `${val}` }
        });
        chart
          .line()
          .position('date*count')
          .color('type')
          .shape('smooth')
          .tooltip('type*count', (type, count) => ({
            name: type,
            value: count + ""
          }))
        chart
          .point()
          .position('date*count')
          .color('type')
          .size(4)
          .shape('circle')
          .style({
            stroke: '#fff',
            lineWidth: 1
          });
        chart.render();
      },

    }
  }
</script>
