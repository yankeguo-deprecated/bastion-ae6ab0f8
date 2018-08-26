<template>
  <b-row class="mt-3">
    <b-col>
      <b-row>
        <b-col>
          <b-breadcrumb :items="navigationItems"/>
        </b-col>
      </b-row>
      <b-row>
        <b-col md="4" lg="3">
          <b-row>
            <b-col>
              <b-card header="用户信息" header-tag="b">
                <b-form @submit.prevent="onProfileFormSubmit">
                  <b-form-group label="账户名" label-class="text-right" horizontal>
                    <b-form-input :value="user.account" readonly plaintext></b-form-input>
                  </b-form-group>
                  <b-form-group label="账户类型" label-class="text-right" horizontal>
                    <b-form-input :value="user | formatUserStatus" readonly plaintext></b-form-input>
                  </b-form-group>
                  <b-form-group label="昵称" label-class="text-right" horizontal>
                    <b-form-input v-model="userNickname" :disabled="busy"></b-form-input>
                  </b-form-group>
                  <b-form-group label="创建时间" label-class="text-right" horizontal>
                    <b-form-input :value="user.created_at | formatUnixEpoch" readonly plaintext></b-form-input>
                  </b-form-group>
                  <b-button type="submit" class="btn-block" :disabled="busy" variant="primary"><i class="fa fa-upload"
                                                                                                  aria-hidden="true"></i> 修改昵称</b-button>
                </b-form>
              </b-card>
            </b-col>
          </b-row>
          <b-row class="mt-3">
            <b-col>
              <b-card v-if="user.account !== currentUser.account" header="操作" header-tag="b">
                  <b-button :disabled="busy" class="btn-block" variant="danger"
                          v-if="user.is_admin && user.account !== accountToDowngrade"
                          @click="onDowngradeClick(user.account)"><i class="fa fa-level-down" aria-hidden="true"></i>
                    降级管理员
                  </b-button>
                  <b-button :disabled="busy" class="btn-block" variant="danger"
                          v-if="user.is_admin && user.account === accountToDowngrade"
                          @click="onDowngradeConfirmClick(user.account)"><i class="fa fa-level-down"
                                                                            aria-hidden="true"></i> 确认降级管理员
                  </b-button>
                  <b-button :disabled="busy" class="btn-block" variant="success"
                          v-if="!user.is_admin && user.account !== accountToUpgrade"
                          @click="onUpgradeClick(user.account)"><i class="fa fa-level-up" aria-hidden="true"></i> 升级管理员
                  </b-button>
                  <b-button :disabled="busy" class="btn-block" variant="success"
                          v-if="!user.is_admin && user.account === accountToUpgrade"
                          @click="onUpgradeConfirmClick(user.account)"><i class="fa fa-level-up"
                                                                          aria-hidden="true"></i> 确认升级管理员
                  </b-button>
                  <b-button :disabled="busy" class="btn-block" variant="success"
                          v-if="user.is_blocked && user.account !== accountToUnblock"
                          @click="onUnblockClick(user.account)"><i class="fa fa-check-circle-o" aria-hidden="true"></i>
                    解封用户
                  </b-button>
                  <b-button :disabled="busy" class="btn-block" variant="success"
                          v-if="user.is_blocked && user.account === accountToUnblock"
                          @click="onUnblockConfirmClick(user.account)"><i class="fa fa-check-circle-o"
                                                                          aria-hidden="true"></i> 确认解封用户
                  </b-button>
                  <b-button :disabled="busy" class="btn-block" variant="danger"
                          v-if="!user.is_blocked && user.account !== accountToBlock"
                          @click="onBlockClick(user.account)"><i class="fa fa-ban" aria-hidden="true"></i> 封禁用户
                  </b-button>
                  <b-button :disabled="busy" class="btn-block" variant="danger"
                          v-if="!user.is_blocked && user.account === accountToBlock"
                          @click="onBlockConfirmClick(user.account)"><i class="fa fa-ban" aria-hidden="true"></i>
                    确认封禁用户
                  </b-button>
              </b-card>
            </b-col>
          </b-row>
        </b-col>
        <b-col md="8" lg="9">
          <b-row>
            <b-col>
              <b-card header="管理权限" header-tag="b" no-body>
                <b-card-body>
                  <b-form inline @submit.prevent="onSubmit">
                    <b-form-select v-model="form.user_mode" :options="user_modes" class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0"></b-form-select>
                    <b-input class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0" v-if="form.user_mode == 'console'" v-model="form.user" placeholder="Linux 用户"/>
                    <i v-if="form.user_mode == 'console'" class="fa fa-at" aria-hidden="true"></i>
                    <b-input class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0" v-model="form.hostname_pattern" placeholder="主机名，允许通配符 *"/>
                    <span>,</span>
                    <b-input v-if="form.expires_mode != 'n'" class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0" v-model="form.expires_in" type="number"/>
                    <b-form-select v-model="form.expires_mode" :options="expire_modes" class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0"></b-form-select>
                    <b-button type="submit" variant="success"><i class="fa fa-pencil-square-o" aria-hidden="true"></i>
                      添加/更新
                    </b-button>
                  </b-form>
                </b-card-body>
                <b-table :items="grants" :fields="fields" class="mb-0" empty-text="无" :show-empty="true">
                  <template slot="type" slot-scope="data">
                    <span v-if="data.item.user === '__tunnel__'"><i class="fa fa-link" aria-hidden="true"></i> 建立隧道</span>
                    <span v-if="data.item.user !== '__tunnel__'"><i class="fa fa-sign-in" aria-hidden="true"></i> 登录用户</span>
                  </template>
                  <template slot="user" slot-scope="data">
                    <span v-if="data.item.user === '__tunnel__'">-</span>
                    <span v-if="data.item.user !== '__tunnel__'">{{data.item.user}}</span>
                  </template>
                  <template slot="created_at" slot-scope="data">
                    {{data.item.created_at | formatUnixEpoch}}
                  </template>
                  <template slot="expired_at" slot-scope="data">
                    {{data.item.expired_at | formatUnixEpochExpired}}
                  </template>
                  <template slot="action" slot-scope="data">
                    <b-link href="#" class="text-danger"
                            v-if="data.item.user != grantToDelete.user || data.item.hostname_pattern != grantToDelete.hostname_pattern"
                            @click="onDeleteClick(data.item)"><i class="fa fa-trash" aria-hidden="true"></i> 删除
                    </b-link>
                    <b-link href="#" class="text-danger"
                            v-if="data.item.user == grantToDelete.user && data.item.hostname_pattern == grantToDelete.hostname_pattern"
                            @click="onDeleteConfirmClick(data.item)"><i class="fa fa-trash" aria-hidden="true"></i> 确认删除
                    </b-link>
                  </template>
                </b-table>
              </b-card>
            </b-col>
          </b-row>
        </b-col>
      </b-row>
    </b-col>
  </b-row>
</template>

<script>
/* eslint-disable camelcase */
import {mapState} from 'vuex'

export default {
  name: 'UserDetail',
  computed: {
    ...mapState(['currentUser'])
  },
  data () {
    return {
      navigationItems: [
        {
          text: '全部用户',
          to: '/users'
        },
        {
          text: '用户: ' + this.$route.params.account,
          active: true
        }
      ],
      fields: [
        {
          key: 'type',
          label: '类型',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'user',
          sortable: true,
          label: 'Linux 用户',
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'hostname_pattern',
          sortable: true,
          label: '主机名',
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'created_at',
          sortable: true,
          label: '创建时间',
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'expired_at',
          sortable: true,
          label: '过期时间',
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
      user: {
        account: this.$route.params.account,
        nickname: this.$route.params.account,
        is_admin: false,
        is_blocked: false,
        created_at: 0,
        updated_at: 0
      },
      userNickname: '',
      busy: false,
      grants: [],
      form: {
        user_mode: 'console',
        user: 'root',
        hostname_pattern: '',
        expires_in: 1,
        expires_mode: 'h'
      },
      user_modes: [
        {
          value: 'tunnel',
          text: '建立隧道'
        },
        {
          value: 'console',
          text: '登录用户'
        }
      ],
      expire_modes: [
        {
          value: 'h',
          text: '小时后过期'
        },
        {
          value: 'd',
          text: '天后过期'
        },
        {
          value: 'n',
          text: '永不过期'
        }
      ],
      grantToDelete: {},
      accountToBlock: '',
      accountToUnblock: '',
      accountToUpgrade: '',
      accountToDowngrade: ''
    }
  },
  mounted () {
    this.fetchUser()
    this.fetchUserGrants()
  },
  methods: {
    fetchUser () {
      this.busy = true
      this.$apiGetUser(this.$route.params.account).then(res => {
        this.busy = false
        this.user = res.body.user
        this.userNickname = this.user.nickname
      }, (res) => {
        this.busy = false
      })
    },
    fetchUserGrants () {
      this.$apiGetUserGrants(this.$route.params.account).then(res => {
        this.grants = res.body.grants || []
      })
    },
    onSubmit () {
      let expires_in = this.form.expires_in
      switch (this.form.expires_mode) {
        case 'h': {
          expires_in *= 3600
          break
        }
        case 'd': {
          expires_in *= 3600 * 24
          break
        }
        case 'n': {
          expires_in = 0
          break
        }
        default: {
          throw Error('not possible value')
        }
      }
      let user = this.form.user
      if (this.form.user_mode === 'tunnel') {
        user = '__tunnel__'
      }
      this.$apiCreateGrant({
        account: this.$route.params.account,
        hostname_pattern: this.form.hostname_pattern,
        user,
        expires_in
      }).then(res => {
        this.fetchUserGrants()
      })
    },
    onProfileFormSubmit () {
      this.busy = true
      this.$apiUpdateUserNickname({account: this.user.account, nickname: this.userNickname}).then((res) => {
        this.user = res.body.user
        this.userNickname = this.user.nickname
        this.busy = false
      }, (res) => {
        this.busy = false
      })
    },
    onDeleteClick (grant) {
      this.grantToDelete = grant
    },
    onDeleteConfirmClick (grant) {
      this.$apiDestroyGrant(this.grantToDelete).then(() => {
        this.fetchUserGrants()
      })
    },
    clearActionStates () {
      this.accountToBlock = ''
      this.accountToUnblock = ''
      this.accountToUpgrade = ''
      this.accountToDowngrade = ''
    },
    onBlockClick (account) {
      this.clearActionStates()
      this.accountToBlock = account
    },
    onBlockConfirmClick (account) {
      this.clearActionStates()
      this.busy = true
      this.$apiUpdateUserIsBlocked({account, is_blocked: true}).then((res) => {
        this.user = res.body.user
        this.busy = false
      }, (res) => {
        this.busy = false
      })
    },
    onUnblockClick (account) {
      this.clearActionStates()
      this.accountToUnblock = account
    },
    onUnblockConfirmClick (account) {
      this.clearActionStates()
      this.busy = true
      this.$apiUpdateUserIsBlocked({account, is_blocked: false}).then((res) => {
        this.user = res.body.user
        this.busy = false
      }, (res) => {
        this.busy = false
      })
    },
    onUpgradeClick (account) {
      this.clearActionStates()
      this.accountToUpgrade = account
    },
    onUpgradeConfirmClick (account) {
      this.clearActionStates()
      this.busy = true
      this.$apiUpdateUserIsAdmin({account, is_admin: true}).then((res) => {
        this.user = res.body.user
        this.busy = false
      }, (res) => {
        this.busy = false
      })
    },
    onDowngradeClick (account) {
      this.clearActionStates()
      this.accountToDowngrade = account
    },
    onDowngradeConfirmClick (account) {
      this.clearActionStates()
      this.busy = true
      this.$apiUpdateUserIsAdmin({account, is_admin: false}).then((res) => {
        this.user = res.body.user
        this.busy = false
      }, (res) => {
        this.busy = false
      })
    }
  }
}
</script>

<style scoped></style>
