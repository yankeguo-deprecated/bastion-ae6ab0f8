import Vue from 'vue'

// install BootstrapVue
import BootstrapVue from 'bootstrap-vue'
import 'bootstrap/dist/css/bootstrap.css'
import 'bootstrap-vue/dist/bootstrap-vue.css'

import router from '@/router'
import store from '@/store'
import App from '@/App'
import '@/api'

import Notifications from 'vue-notification'

Vue.use(BootstrapVue)
Vue.use(Notifications)

Vue.config.productionTip = false

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: {App},
  template: '<App/>'
})
