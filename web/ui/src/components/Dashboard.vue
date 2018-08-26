<template>
  <b-row class="mt-3">
    <b-col md="4" lg="2">
      <b-list-group>
        <b-list-group-item :active="index === 0" @click="index = 0" href="#">
          <i class="fa fa-sign-in" aria-hidden="true"></i> 登录服务器终端
        </b-list-group-item>
        <b-list-group-item :active="index === 1" @click="index = 1" href="#">
          <i class="fa fa-link" aria-hidden="true"></i> 建立 TCP 隧道
        </b-list-group-item>
      </b-list-group>
    </b-col>
    <b-col md="8" lg="10">
      <b-row v-if="index === 0">
        <b-col>
          <b-card header="登录沙箱环境" header-tag="b">
            <p>使用如下命令登录沙箱环境（不需要指定用户名）</p>
            <p class="mb-0"><code>ssh {{ssh_domain}}</code></p>
          </b-card>
        </b-col>
      </b-row>
      <b-row class="mt-3" v-if="index === 0">
        <b-col>
          <b-card no-body header="从沙箱环境登录服务器" header-tag="b">
            <b-card-body>
              <p class="mb-0">当前有权限登录的服务器:</p>
            </b-card-body>
            <b-table :items="grantItems" :fields="fields" class="mb-0" :show-empty="true" empty-text="无">
              <template slot="command" slot-scope="data">
                <code>ssh {{data.item.hostname}}-{{data.item.user}}</code>
              </template>
            </b-table>
          </b-card>
        </b-col>
      </b-row>
      <b-row v-if="index === 1">
        <b-col>
          <b-card header="使用 SSH 建立 TCP 隧道" header-tag="b" class="mb-2" no-body>
            <b-card-body>
              <p>使用堡垒机可以建立本地端口到远程服务器的端口的 TCP 隧道（无需进入沙箱环境）</p>
              <b-card bg-variant="light">
                <p>假设需要在本地访问 <code>example.db.01</code> 的 3306 端口，运行如下命令</p>
                <p><code>ssh -N -L 4306:example.db.01:3306 {{ssh_domain}}</code></p>
                <p class="mb-0">该命令会监听本地的 4306 端口（可以按需要选择任意端口），访问该端口等同于访问 <code>example.db.01</code> 的 3306 端口</p>
              </b-card>
              <p class="mt-3 mb-0">当前有权限建立TCP隧道的服务器：</p>
            </b-card-body>
            <b-table :items="grantTunnels" :fields="fieldsTunnels" class="mb-0" :show-empty="true" empty-text="无">
              <template slot="command" slot-scope="data">
                <code>ssh -N -L $LOCAL_PORT:{{data.item.hostname}}:$REMOTE_PORT {{ssh_domain}}</code>
              </template>
            </b-table>
          </b-card>
        </b-col>
      </b-row>
    </b-col>
  </b-row>
</template>

<script>
import {mapState} from 'vuex'

export default {
  name: 'Dashboard',
  data () {
    return {
      index: 0,
      fields: [
        {
          key: 'user',
          label: 'Linux 用户',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'hostname',
          label: '主机名',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'command',
          label: '登录终端命令',
          thClass: 'text-center'
        }
      ],
      fieldsTunnels: [
        {
          key: 'hostname',
          label: '主机名',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'command',
          label: '建立 TCP 隧道命令',
          thClass: 'text-center'
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
  }
}
</script>

<style scoped></style>
