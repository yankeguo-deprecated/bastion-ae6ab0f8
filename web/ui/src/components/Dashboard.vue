<template>
  <b-row class="mt-3">
    <b-col>
      <b-tabs pills vertical nav-wrapper-class="dashboard-sidebar" @input="onTabSwitched">
        <b-tab active>
          <template slot="title">
            <i class="fa fa-sign-in" aria-hidden="true"></i> 登录服务器终端
          </template>
          <b-row>
            <b-col>
                <b-card title="1. 登录沙箱环境" title-tag="h5">
                  <p>使用如下命令登录沙箱环境（不需要指定用户名）</p>
                  <p><code>ssh {{ssh_domain}}</code></p>
                </b-card>
            </b-col>
          </b-row>
          <b-row class="mt-3">
            <b-col>
              <b-card title="2. 从沙箱环境登录服务器" title-tag="h5">
                <p>当前有权限连接的服务器:</p>
                <b-table v-if="grantItems.length > 0" :items="grantItems" :fields="fields">
                  <template slot="command" slot-scope="data">
                    <code>ssh {{data.item.hostname}}-{{data.item.user}}</code>
                  </template>
                </b-table>
                <div v-if="grantItems.length === 0" class="text-center">
                  <hr/>
                  <span class="text-muted">无</span>
                  <hr/>
                </div>
              </b-card>
            </b-col>
          </b-row>
        </b-tab>
        <b-tab>
          <template slot="title">
            <i class="fa fa-link" aria-hidden="true"></i> 建立 TCP 隧道
          </template>
          <b-card title="1. 使用 SSH 建立 TCP 隧道" title-tag="h5" class="mb-2">
            <p>使用堡垒机可以建立本地端口到远程服务器的端口的 TCP 隧道（无需进入沙箱环境）</p>
            <b-card bg-variant="light">
              <p>假设需要在本地访问 <code>example.db.01</code> 的 3306 端口，运行如下命令</p>
              <p><code>ssh -N -L 4306:example.db.01:3306 {{ssh_domain}}</code></p>
              <p class="mb-0">该命令会监听本地的 4306 端口（可以按需要选择任意端口），访问该端口等同于访问 <code>example.db.01</code> 的 3306 端口</p>
            </b-card>
            <p class="mt-3">当前有权限建立TCP隧道的服务器：</p>
            <b-table v-if="grantTunnels.length > 0" :items="grantTunnels" :fields="fieldsTunnels">
              <template slot="command" slot-scope="data">
                <code>ssh -N -L $LOCAL_PORT:{{data.item.hostname}}:$REMOTE_PORT {{ssh_domain}}</code>
              </template>
            </b-table>
            <div v-if="grantTunnels.length === 0" class="text-center">
              <hr/>
              <span class="text-muted">无</span>
              <hr/>
            </div>
          </b-card>
        </b-tab>
      </b-tabs>
    </b-col>
  </b-row>
</template>

<script>
import {mapState} from 'vuex'

export default {
  name: 'Dashboard',
  data () {
    return {
      fields: [
        {
          key: 'user',
          label: 'Linux 用户',
          sortable: true
        },
        {
          key: 'hostname',
          label: '主机名',
          sortable: true
        },
        {
          key: 'command',
          label: '连接命令'
        }
      ],
      fieldsTunnels: [
        {
          key: 'hostname',
          label: '主机名',
          sortable: true
        },
        {
          key: 'command',
          label: '建立隧道命令'
        }
      ]
    }
  },
  mounted () {
    this.$apiGetCurrentUser()
    this.$apiGetCurrentUserGrantItems()
  },
  computed: {
    ...mapState(['ssh_domain', 'grantItems', 'grantTunnels'])
  },
  methods: {
    onTabSwitched (index) {
    }
  }
}
</script>

<style>
  .dashboard-sidebar {
    width: 16rem;
  }
</style>
