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
    nodes: [],
    keys: []
  },
  getters: {
    isLoggedIn: state => !!state.currentToken && !!state.currentUser,
    isLoggedInAsAdmin: state =>
      !!state.currentToken && !!state.currentUser && state.currentUser.is_admin
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
    },
    setKeys (state, keys) {
      state.keys = keys || []
    },
    setNodes (state, nodes) {
      state.nodes = nodes || []
    },
    setUsers (state, users) {
      state.users = users || []
    }
  },
  plugins: [createPersistedState()]
})
