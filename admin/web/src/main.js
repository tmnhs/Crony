import 'babel-polyfill'
import Vue from 'vue'
import dayjs from 'dayjs'
import ElementUI from 'element-ui'
import reminder from './components/base/Reminder'
import App from './App'
import router from './router'
import store from './store'
import api from './api'
import './assets/styles/app.scss'
import './assets/icons'
import './components'
import './directive'
import './filters'
import i18n from './assets/lang'

Vue.use(ElementUI, {
	size: store.getters.size,
})

Object.defineProperty(Vue.prototype, '$api', {
	value: api,
})
Object.defineProperty(Vue.prototype, '$dayjs', {
	value: dayjs,
})
Object.defineProperty(Vue.prototype, '$reminder', {
	value: reminder,
})

Vue.config.productionTip = false

new Vue({
	el: '#app',
	i18n,
	router: router,
	store: store,
	render: h => h(App),
})
