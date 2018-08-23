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
          <b-card>
            <b-form>
              <b-form-group label="账户名:" label-class="text-right" horizontal>
                <b-form-input :value="user.account" readonly plaintext></b-form-input>
              </b-form-group>
              <b-form-group label="账户类型:" label-class="text-right" horizontal>
                <b-form-input :value="user | formatUserStatus" readonly plaintext></b-form-input>
              </b-form-group>
              <b-form-group label="昵称:" label-class="text-right" horizontal>
                <b-form-input :value="user.nickname" readonly plaintext></b-form-input>
              </b-form-group>
              <b-form-group label="创建时间:" label-class="text-right" horizontal>
                <b-form-input :value="user.created_at | formatUnixEpoch" readonly plaintext></b-form-input>
              </b-form-group>
            </b-form>
          </b-card>
        </b-col>
        <b-col md="8" lg="9">
          <b-row>
            <b-col>
              <b-card>
                <b-form inline @submit="onSubmit">
                  <b-form-select v-model="form.user_mode" :options="user_modes"
                                 class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0"></b-form-select>
                  <b-input class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0" v-if="form.user_mode == 'console'" v-model="form.user" placeholder="Linux 用户"/>
                  <i v-if="form.user_mode == 'console'" class="fa fa-at" aria-hidden="true"></i>
                  <b-input class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0" v-model="form.hostname_pattern"
                           placeholder="主机名，允许通配符 *"/>
                  <span>,</span>
                  <b-input v-if="form.expires_mode != 'n'" class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0"
                           v-model="form.expires_in" type="number"/>
                  <b-form-select v-model="form.expires_mode" :options="expire_modes"
                                 class="ml-sm-2 mb-2 mr-sm-2 mb-sm-0"></b-form-select>
                  <b-button type="submit" variant="success"><i class="fa fa-pencil-square-o" aria-hidden="true"></i>
                    添加/更新
                  </b-button>
                </b-form>
                <b-table :items="grants" :fields="fields" class="mt-3">
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

export default {
  name: 'UserDetail',
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
          key: 'expired_at',
          sortable: true,
          label: '过期时间',
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
      grantToDelete: {}
    }
  },
  mounted () {
    this.fetchUser()
    this.fetchUserGrants()
  },
  methods: {
    fetchUser () {
      this.$apiGetUser(this.$route.params.account).then(res => {
        this.user = res.body.user
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
    onDeleteClick (grant) {
      this.grantToDelete = grant
    },
    onDeleteConfirmClick (grant) {
      this.$apiDestroyGrant(this.grantToDelete).then(() => {
        this.fetchUserGrants()
      })
    }
  }
}
</script>

<style scoped>
</style>
