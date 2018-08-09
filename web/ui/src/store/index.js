import Vue from 'vue'
import Vuex from 'vuex'
import createPersistedState from 'vuex-persistedstate'

Vue.use(Vuex)

export default new Vuex.Store({
  state: {
    currentToken: null,
    currentUser: null,
    tokens: [],
    users: [],
    grantItems: [],
    nodes: []
  },
  getters: {
    isLoggedIn: state => (!!state.currentToken) && (!!state.currentUser),
    isLoggedInAsAdmin: state => (!!state.currentToken) && (!!state.currentUser) && state.currentUser.is_admin,
    userType: state => {
      let types = []
      if (state.currentUser) {
        if (state.currentUser.is_admin) {
          types.push('管理员')
        } else {
          types.push('普通用户')
        }
        if (state.currentUser.is_blocked) {
          types.push('被封禁')
        }
      }
      return types.join(',')
    }
  },
  mutations: {
    setCurrentToken (state, token) {
      state.currentToken = token
    },
    setCurrentUser (state, user) {
      state.currentUser = user
    },
    setGrantItems (state, gis) {
      state.grantItems = gis || []
    },
    setTokens (state, tokens) {
      state.tokens = tokens || []
    }
  },
  plugins: [createPersistedState()]
})
