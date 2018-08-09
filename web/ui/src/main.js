import Vue from 'vue'
import moment from 'moment'
import UAParser from 'ua-parser-js'

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

Vue.filter('formatUnixEpoch', function (value) {
  if (!value) {
    return '-'
  }
  return moment(value * 1000).format('YYYY-MM-DD HH:mm:ss')
})

Vue.filter('formatUserAgent', function (ua) {
  let { browser, os } = UAParser(ua)
  return `${browser.name} ${browser.version} (${os.name} ${os.version})`
})

Vue.filter('formatUserStatus', function (u) {
  if (!u) {
    return ''
  } else {
    let types = []
    if (u.is_admin) {
      types.push('管理员')
    }
    if (u.is_blocked) {
      types.push('已封禁')
    }
    return types.join(', ')
  }
})

/* eslint-disable no-new */
new Vue({
  el: '#app',
  router,
  store,
  components: { App },
  template: '<App/>'
})
