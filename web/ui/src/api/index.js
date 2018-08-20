/* eslint-disable camelcase */
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
      return res => {
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
      return this.$http
        .post('/api/tokens/create', data, {emulateJSON: true})
        .then(res => {
          store.commit('setCurrentUser', res.body.user)
          store.commit('setCurrentToken', res.body.token)
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiLogout = function () {
      return this.$http
        .post(
          '/api/tokens/destroy',
          {id: store.state.currentToken.id},
          {emulateJSON: true}
        )
        .then(res => {
          store.commit('setCurrentUser', null)
          store.commit('setCurrentToken', null)
          this.$router.push('/login')
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiGetCurrentUser = function () {
      return this.$http.get('/api/users/current').then(res => {
        store.commit('setCurrentUser', res.body.user)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiGetCurrentUserGrantItems = function () {
      return this.$http.get('/api/users/current/grant_items').then(res => {
        store.commit('setGrantItems', res.body.grant_items)
        store.commit('setSSHDomain', res.body.ssh_domain)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiUpdateCurrentUserNickname = function (nickname) {
      return this.$http
        .post(
          '/api/users/current/update_nickname',
          {nickname},
          {emulateJSON: true}
        )
        .then(res => {
          store.commit('setCurrentUser', res.body.user)
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: '昵称修改成功'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiUpdatePassword = function ({oldPassword, newPassword}) {
      return this.$http
        .post(
          '/api/users/current/update_password',
          {oldPassword, newPassword},
          {emulateJSON: true}
        )
        .then(res => {
          store.commit('setCurrentUser', res.body.user)
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: '密码修改成功'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiListTokens = function () {
      return this.$http.get('/api/tokens').then(res => {
        store.commit('setTokens', res.body.tokens)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiDeleteToken = function (id) {
      return this.$http
        .post('/api/tokens/destroy', {id}, {emulateJSON: true})
        .then(res => {
          this.$apiListTokens()
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: '访问令牌已删除'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiListKeys = function () {
      return this.$http.get('/api/users/current/keys').then(res => {
        store.commit('setKeys', res.body.keys)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiCreateKey = function ({name, publicKey}) {
      return this.$http
        .post(
          '/api/users/current/keys/create',
          {name, publicKey},
          {emulateJSON: true}
        )
        .then(res => {
          this.$apiListKeys()
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: 'SSH 公钥已添加'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiDestroyKey = function (fingerprint) {
      return this.$http
        .post('/api/keys/destroy', {fingerprint}, {emulateJSON: true})
        .then(res => {
          this.$apiListKeys()
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: 'SSH 公钥已移除'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiListNodes = function () {
      return this.$http.get('/api/nodes').then(res => {
        store.commit('setNodes', res.body.nodes)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiCreateNode = function ({hostname, address}) {
      return this.$http
        .post('/api/nodes/create', {hostname, address}, {emulateJSON: true})
        .then(res => {
          this.$apiListNodes()
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: '服务器已添加/更新'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiDestroyNode = function (hostname) {
      return this.$http
        .post('/api/nodes/destroy', {hostname}, {emulateJSON: true})
        .then(res => {
          this.$apiListNodes()
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: '服务器已移除'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiListUsers = function () {
      return this.$http.get('/api/users').then(res => {
        store.commit('setUsers', res.body.users)
        return res
      }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiCreateUser = function ({account, nickname, password}) {
      return this.$http
        .post(
          '/api/users/create',
          {account, nickname, password},
          {emulateJSON: true}
        )
        .then(res => {
          this.$apiListUsers()
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: '用户已添加'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiUpdateUserIsAdmin = function ({account, is_admin}) {
      return this.$http
        .post(
          '/api/users/update_is_admin',
          {account, is_admin},
          {emulateJSON: true}
        )
        .then(res => {
          this.$apiListUsers()
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: is_admin ? '用户已提升为管理员' : '用户已降级管理员'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiUpdateUserIsBlocked = function ({account, is_blocked}) {
      return this.$http
        .post(
          '/api/users/update_is_blocked',
          {account, is_blocked},
          {emulateJSON: true}
        )
        .then(res => {
          this.$apiListUsers()
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: is_blocked ? '用户已解封' : '用户已封禁'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiGetUser = function (account) {
      return this.$http
        .get(`/api/users/${account}`)
        .then(null, this.$apiErrorCallback())
    }
    Vue.prototype.$apiGetUserGrants = function (account) {
      return this.$http
        .get(`/api/users/${account}/grants`)
        .then(null, this.$apiErrorCallback())
    }
    Vue.prototype.$apiCreateGrant = function ({
      account,
      hostname_pattern,
      user,
      expires_in
    }) {
      return this.$http
        .post(
          `/api/users/${account}/grants/create`,
          {hostname_pattern, user, expires_in},
          {emulateJSON: true}
        )
        .then(res => {
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: '授权已添加/更新'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiDestroyGrant = function ({
      account,
      hostname_pattern,
      user
    }) {
      return this.$http
        .post(
          `/api/users/${account}/grants/destroy`,
          {hostname_pattern, user},
          {emulateJSON: true}
        )
        .then(res => {
          this.$notify({
            type: 'success',
            title: '操作成功',
            text: '授权已移除'
          })
          return res
        }, this.$apiErrorCallback())
    }
    Vue.prototype.$apiListSessions = function ({skip, limit}) {
      return this.$http.get('/api/sessions', {params: {skip, limit}}).then(null, this.$apiErrorCallback())
    }
    Vue.prototype.$apiDownloadReplay = function (id) {
      return this.$http.get(`/api/replays/${id}/download`).then(null, this.$apiErrorCallback())
    }
  }
}

Vue.use(API)
