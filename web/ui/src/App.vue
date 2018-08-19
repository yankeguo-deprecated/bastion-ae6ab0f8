<template>
  <div id="app">
    <notifications position="top center"/>
    <b-navbar fixed :sticky="true" toggleable="md" type="dark" variant="primary">

      <b-navbar-toggle target="nav_collapse"></b-navbar-toggle>

      <b-navbar-brand to="/">Bastion</b-navbar-brand>

      <b-collapse is-nav id="nav_collapse">

        <b-navbar-nav>
          <b-nav-item v-if="isLoggedIn" to="/dashboard"><i class="fa fa-tachometer" aria-hidden="true"></i> 工作台
          </b-nav-item>
          <b-nav-item v-if="isLoggedInAsAdmin" to="/servers"><i class="fa fa-server" aria-hidden="true"></i> 服务器列表
          </b-nav-item>
          <b-nav-item v-if="isLoggedInAsAdmin" to="/users"><i class="fa fa-user-circle-o" aria-hidden="true"></i>
            用户列表
          </b-nav-item>
          <b-nav-item v-if="isLoggedInAsAdmin" to="/sessions"><i class="fa fa-list-alt" aria-hidden="true"></i> 操作记录
          </b-nav-item>
        </b-navbar-nav>

        <b-navbar-nav class="ml-auto">
          <b-nav-item-dropdown right v-if="isLoggedIn">
            <template slot="button-content">
              <i class="fa fa-user-circle" aria-hidden="true"></i> {{currentUser.nickname}} ({{currentUser.account}})
            </template>
            <b-dropdown-item to="/settings/profile"><i class="fa fa-user-circle" aria-hidden="true"></i> 个人设置
            </b-dropdown-item>
            <b-dropdown-divider></b-dropdown-divider>
            <b-dropdown-item @click="onLogoutClick"><i class="fa fa-sign-out" aria-hidden="true"></i> 退出登录
            </b-dropdown-item>
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
    return {}
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
    font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB",
    "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
  }

  td.action-cell {
    width: 8rem;
    text-align: center;
  }

  td.action-cell-wide {
    width: 16rem;
    text-align: center;
  }
</style>
