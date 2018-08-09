<template>
  <div id="app">
    <notifications position="top center"/>
    <b-navbar fixed :sticky="true" toggleable="md" type="dark" variant="primary">

      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>

      <b-navbar-brand to="/">Bastion</b-navbar-brand>

      <b-collapse is-nav id="nav_collapse">

        <b-navbar-nav>
          <b-nav-item v-if="isLoggedIn" to="/dashboard">工作台</b-nav-item>
          <b-nav-item v-if="isLoggedInAsAdmin" to="/servers">服务器列表</b-nav-item>
          <b-nav-item v-if="isLoggedInAsAdmin" to="/users">用户列表</b-nav-item>
          <b-nav-item v-if="isLoggedInAsAdmin" to="/sessions">操作记录</b-nav-item>
        </b-navbar-nav>

        <b-navbar-nav class="ml-auto">
          <b-nav-item-dropdown right v-if="isLoggedIn">
            <template slot="button-content">
              <em>当前用户: {{currentUser.nickname}} ({{currentUser.account}})</em>
            </template>
            <b-dropdown-item to="/settings/profile">个人设置</b-dropdown-item>
            <b-dropdown-item-divider></b-dropdown-item-divider>
            <b-dropdown-item @click="onLogoutClick">退出登录</b-dropdown-item>
          </b-nav-item-dropdown>
        </b-navbar-nav>
      </b-collapse>
    </b-navbar>
    <b-container fluid>
      <router-view/>
    </b-container>
  </div>
</template>

<script>
import {mapGetters, mapState} from 'vuex'

export default {
  name: 'App',
  data () {
    return {
    }
  },
  computed: {
    ...mapState(['currentToken', 'currentUser']),
    ...mapGetters(['isLoggedIn', 'isLoggedInAsAdmin'])
  },
  methods: {
    onLogoutClick () {
      this.$apiLogout()
    }
  }
}
</script>

<style>
  body {
    font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  a.destroy-link {
    color: #ff4757;
  }
</style>
