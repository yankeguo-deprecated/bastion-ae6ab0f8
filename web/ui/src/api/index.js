import Vue from 'vue'
import VueResource from 'vue-resource'

import store from '@/store'

Vue.use(VueResource)

Vue.http.interceptors.push(function (request) {
  if (store.state.currentToken) {
    request.headers.set('X-Bastion-Token', store.state.currentToken.token)
  }
  return function (response) {
    if (response.headers.get('X-Bastion-Action') === 'clear-token') {
      store.commit('setCurrentToken', null)
      store.commit('setCurrentUser', null)
    }
  }
})

let API = {
  install (Vue, options) {
    Vue.prototype.$apiErrorCallback = function () {
      return (res) => {
        console.log(res)
        this.$notify({
          type: 'error',
          title: 'API Error',
          text: res.bodyText,
          duration: 2000
        })
        return res
      }
    }
    Vue.prototype.$apiLogin = function (data) {
      return this.$http.post('/api/tokens/create', data, {emulateJSON: true}).then((res) => {
        store.commit('setCurrentUser', res.body.user)
        store.commit('setCurrentToken', res.body.token)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiGetCurrentUser = function () {
      return this.$http.get('/api/users/current').then((res) => {
        store.commit('setCurrentUser', res.body.user)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiGetCurrentUserGrantItems = function () {
      return this.$http.get('/api/users/current/grant_items').then((res) => {
        store.commit('setGrantItems', res.body.grant_items)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiUpdateCurrentUserNickname = function (nickname) {
      return this.$http.post('/api/users/current/update_nickname', {nickname}, {emulateJSON: true}).then((res) => {
        store.commit('setCurrentUser', res.body.user)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiListTokens = function () {
      return this.$http.get('/api/tokens/list').then((res) => {
        store.commit('setTokens', res.body.tokens)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiDeleteToken = function (id) {
      return this.$http.post('/api/tokens/destroy', {id}, {emulateJSON: true})
    }
  }
}

Vue.use(API)
