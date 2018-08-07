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
    isLoggedInAsAdmin: state => (!!state.currentToken) && (!!state.currentUser) && state.currentUser.is_admin
  },
  mutations: {
    updateCurrentToken (state, token) {
      state.currentToken = token
    },
    updateCurrentUser (state, user) {
      state.currentUser = user
    }
  },
  plugins: [createPersistedState()]
})
