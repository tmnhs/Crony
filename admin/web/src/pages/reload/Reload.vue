<template></template>

<script>
// 用来刷新路由
export default {
	name:"Reload",
	beforeRouteEnter(to, from, next) {
		next(vm => {
			vm.$router.replace(from.path)
		})
	},
}
</script>
