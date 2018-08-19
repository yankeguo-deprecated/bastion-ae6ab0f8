<template>
  <b-row class="mt-3">
    <b-col md="4" lg="3">
      <b-card title="1. 连接沙箱环境" title-tag="h5">
        <p>使用如下命令连接沙箱环境（不需要指定用户名）</p>
        <p><code>ssh {{ssh_domain}}</code></p>
      </b-card>
    </b-col>
    <b-col md="8" lg="9">
      <b-card title="2. 从沙箱环境连接服务器" title-tag="h5">
        <p>当前有权限连接的服务器:</p>
        <b-table :items="grantItems" :fields="fields">
          <template slot="command" slot-scope="data">
            <code>ssh {{data.item.hostname}}-{{data.item.user}}</code>
          </template>
        </b-table>
      </b-card>
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
      ]
    }
  },
  mounted () {
    this.$apiGetCurrentUser()
    this.$apiGetCurrentUserGrantItems()
  },
  computed: {
    ...mapState(['ssh_domain', 'grantItems'])
  }
}
</script>

<style scoped>
</style>
