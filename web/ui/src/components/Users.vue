<template>
  <b-row class="mt-4">
    <b-modal id="modal1" title="添加用户" :hide-footer="true">
      <b-form @submit="onCreateSubmit">
        <b-form-group label="用户名:" label-class="text-right" description="仅允许英文数字和'.' '-' '_'" horizontal>
          <b-form-input v-model="form.account" placeholder="请输入用户名" type="text"></b-form-input>
        </b-form-group>
        <b-form-group label="昵称:" label-class="text-right" description="建议使用姓名" horizontal>
          <b-form-input v-model="form.nickname" placeholder="请输入昵称" type="text"></b-form-input>
        </b-form-group>
        <b-form-group label="密码:" label-class="text-right" description="长度必须大于6" horizontal>
          <b-form-input v-model="form.password" placeholder="请输入密码" type="password"></b-form-input>
        </b-form-group>
        <b-form-group label="重复密码:" label-class="text-right" horizontal>
          <b-form-input v-model="form.repPassword" placeholder="重复密码" type="password"></b-form-input>
        </b-form-group>
        <div class="text-right">
          <b-button type="submit" :disabled="busy" variant="success"><i class="fa fa-plus-circle" aria-hidden="true"></i> 添加</b-button>
        </div>
      </b-form>
    </b-modal>
    <b-col>
      <b-row>
        <b-col md="4" lg="3">
          <b-form @reset="onReset" inline>
            <b-input v-model="search" class="mb-2 mr-sm-2 mb-sm-0" placeholder="搜索用户名或昵称"/>
            <b-button type="reset" :disabled="search == ''" variant="outline-danger"><i class="fa fa-ban" aria-hidden="true"></i> 清除</b-button>
          </b-form>
        </b-col>
        <b-col md="8" lg="9" class="text-right">
          <b-btn variant="success" v-b-modal.modal1><i class="fa fa-plus-circle" aria-hidden="true"></i> 添加</b-btn>
        </b-col>
      </b-row>
      <b-row class="mt-4">
        <b-col>
          <b-table striped :items="filteredUsers" :fields="fields">
           <template slot="created_at" slot-scope="data">
              {{data.item.created_at | formatUnixEpoch}}
            </template>
            <template slot="viewed_at" slot-scope="data">
              {{data.item.viewed_at | formatUnixEpoch}}
            </template>
            <template slot="action" slot-scope="data">
              <b-link href="#" class="text-danger" v-if="data.item.account != currentUser.account && data.item.is_admin && data.item.account != accountToDowngrade"
                      @click="onDowngradeClick(data.item.account)"><i class="fa fa-level-down" aria-hidden="true"></i> 降级管理员
              </b-link>
              <b-link href="#" class="text-danger" v-if="data.item.account != currentUser.account && data.item.is_admin && data.item.account == accountToDowngrade"
                      @click="onDowngradeConfirmClick(data.item.account)"><i class="fa fa-level-down" aria-hidden="true"></i> 确认降级管理员
              </b-link>
              <b-link href="#" class="text-success" v-if="data.item.account != currentUser.account && !data.item.is_admin && data.item.account != accountToUpgrade"
                      @click="onUpgradeClick(data.item.account)"><i class="fa fa-level-up" aria-hidden="true"></i> 升级管理员
              </b-link>
              <b-link href="#" class="text-success" v-if="data.item.account != currentUser.account && !data.item.is_admin && data.item.account == accountToUpgrade"
                      @click="onUpgradeConfirmClick(data.item.account)"><i class="fa fa-level-up" aria-hidden="true"></i> 确认升级管理员
              </b-link>
              <span class="text-muted" v-if="data.item.account != currentUser.account">&nbsp;|&nbsp;</span>
              <b-link href="#" class="text-success" v-if="data.item.account != currentUser.account && data.item.is_blocked && data.item.account != accountToUnblock"
                      @click="onUnblockClick(data.item.account)"><i class="fa fa-check-circle-o" aria-hidden="true"></i> 解封用户
              </b-link>
              <b-link href="#" class="text-success" v-if="data.item.account != currentUser.account && data.item.is_blocked && data.item.account == accountToUnblock"
                      @click="onUnblockConfirmClick(data.item.account)"><i class="fa fa-check-circle-o" aria-hidden="true"></i> 确认解封用户
              </b-link>
              <b-link href="#" class="text-danger" v-if="data.item.account != currentUser.account && !data.item.is_blocked && data.item.account != accountToBlock"
                      @click="onBlockClick(data.item.account)"><i class="fa fa-ban" aria-hidden="true"></i> 封禁用户
              </b-link>
              <b-link href="#" class="text-danger" v-if="data.item.account != currentUser.account && !data.item.is_blocked && data.item.account == accountToBlock"
                      @click="onBlockConfirmClick(data.item.account)"><i class="fa fa-ban" aria-hidden="true"></i> 确认封禁用户
              </b-link>

              <span class="text-muted" v-if="data.item.account == currentUser.account">(当前用户)</span>
            </template>
            <template slot="status" slot-scope="data">
              {{data.item | formatUserStatus }}
            </template>
          </b-table>
        </b-col>
      </b-row>
    </b-col>
  </b-row>
</template>

<script>
/* eslint-disable standard/object-curly-even-spacing */

import { mapState } from 'vuex'

export default {
  name: 'Users',
  data () {
    return {
      fields: [
        {
          key: 'account',
          label: '用户名',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'nickname',
          label: '昵称',
          thClass: 'text-center',
          tdClass: 'text-center'
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
        },
        {
          key: 'action',
          label: '    ',
          thClass: 'text-center',
          tdClass: 'action-cell-wide'
        }
      ],
      form: {
        account: '',
        nickname: '',
        password: '',
        repPassword: ''
      },
      search: '',
      busy: false,
      accountToBlock: '',
      accountToUnblock: '',
      accountToUpgrade: '',
      accountToDowngrade: ''
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
      this.$apiCreateUser(this.form)
    },
    onBlockClick (account) {
      this.accountToBlock = account
      this.accountToUnblock = ''
      this.accountToUpgrade = ''
      this.accountToDowngrade = ''
    },
    onBlockConfirmClick (account) {
      this.$apiUpdateUserIsBlocked({ account, is_blocked: true })
    },
    onUnblockClick (account) {
      this.accountToBlock = ''
      this.accountToUnblock = account
      this.accountToUpgrade = ''
      this.accountToDowngrade = ''
    },
    onUnblockConfirmClick (account) {
      this.$apiUpdateUserIsBlocked({ account, is_blocked: false })
    },
    onUpgradeClick (account) {
      this.accountToBlock = ''
      this.accountToUnblock = ''
      this.accountToUpgrade = account
      this.accountToDowngrade = ''
    },
    onUpgradeConfirmClick (account) {
      this.$apiUpdateUserIsAdmin({ account, is_admin: true })
    },
    onDowngradeClick (account) {
      this.accountToBlock = ''
      this.accountToUnblock = ''
      this.accountToUpgrade = ''
      this.accountToDowngrade = account
    },
    onDowngradeConfirmClick (account) {
      this.$apiUpdateUserIsAdmin({ account, is_admin: false })
    }
  }
}
</script>

<style scoped>
</style>
