<template>
  <b-row class="mt-3 mb-3">
    <b-col md="4" lg="3">
      <b-card header="添加用户" header-tag="b">
        <b-form @submit.prevent="onCreateSubmit">
          <b-form-group label="用户名" description="仅允许英文数字和'.' '-' '_'">
            <b-form-input v-model="form.account" placeholder="请输入用户名" type="text"></b-form-input>
          </b-form-group>
          <b-form-group label="昵称" description="建议使用姓名">
            <b-form-input v-model="form.nickname" placeholder="请输入昵称" type="text"></b-form-input>
          </b-form-group>
          <b-form-group label="密码" description="长度必须大于6">
            <b-form-input v-model="form.password" placeholder="请输入密码" type="password"></b-form-input>
          </b-form-group>
          <b-form-group label="重复密码">
            <b-form-input v-model="form.repPassword" placeholder="重复密码" type="password"></b-form-input>
          </b-form-group>
          <b-button type="submit" class="btn-block" :disabled="busy" variant="success">
            <i class="fa fa-plus-circle" aria-hidden="true"></i> 添加
          </b-button>
        </b-form>
      </b-card>
    </b-col>
    <b-col md="8" lg="9">
      <b-row>
        <b-col>
          <b-card no-body header="用户列表" header-tag="b">
            <b-card-body>
              <b-form @reset="onReset" inline>
                <b-input v-model="search" class="mb-2 mr-sm-2 mb-sm-0" placeholder="搜索用户名或昵称"/>
                <b-button type="reset" :disabled="search == ''" variant="outline-danger">
                  <i class="fa fa-ban" aria-hidden="true"></i> 清除
                </b-button>
              </b-form>
            </b-card-body>
            <b-table striped :items="filteredUsers" :fields="fields" class="mb-0">
              <template slot="account" slot-scope="data">
                <b-link :to="{name: 'UserDetail', params: {account: data.item.account}}">{{data.item.nickname}} ({{data.item.account}})</b-link>
                <span class="text-muted pull-right" v-if="data.item.account == currentUser.account">当前用户</span>
              </template>
              <template slot="created_at" slot-scope="data">
                {{data.item.created_at | formatUnixEpoch}}
              </template>
              <template slot="viewed_at" slot-scope="data">
                {{data.item.viewed_at | formatUnixEpoch}}
              </template>
              <template slot="status" slot-scope="data">
                {{data.item | formatUserStatus }}
              </template>
            </b-table>
          </b-card>
        </b-col>
      </b-row>
    </b-col>
  </b-row>
</template>

<script>
/* eslint-disable standard/object-curly-even-spacing */

import {mapState} from 'vuex'

export default {
  name: 'Users',
  data () {
    return {
      fields: [
        {
          key: 'account',
          label: '用户',
          sortable: true,
          thClass: 'text-center'
        },
        {
          key: 'status',
          label: '账户类型',
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'created_at',
          label: '创建日期',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'viewed_at',
          label: '最后使用',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        }
      ],
      form: {
        account: '',
        nickname: '',
        password: '',
        repPassword: ''
      },
      search: '',
      busy: false
    }
  },
  mounted () {
    this.$apiListUsers()
  },
  computed: {
    ...mapState(['users', 'currentUser']),
    filteredUsers () {
      if (this.search && this.search.length > 0) {
        return this.users.filter(n => {
          console.log(n)
          return (
            n.account.includes(this.search) || n.nickname.includes(this.search)
          )
        })
      } else {
        return this.users
      }
    }
  },
  methods: {
    onReset () {
      this.search = ''
    },
    onCreateSubmit () {
      if (this.form.password !== this.form.repPassword) {
        this.$notify({
          type: 'warn',
          title: '输入错误',
          text: '重复密码不正确'
        })
      }
      this.$apiCreateUser(this.form).then((res) => {
        this.form.account = ''
        this.form.nickname = ''
        this.form.password = ''
        this.form.repPassword = ''
      })
    }
  }
}
</script>

<style scoped></style>
