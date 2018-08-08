<template>
  <b-row id="main-row">
    <b-col>
      <b-tabs>
        <b-tab title="基本信息" active class="main-tab">
          <b-row>
            <b-col md="4">
              <b-form @submit="onUpdateSubmit">
                <b-form-group label="账户名" horizontal>
                  <b-form-input :value="currentUser.account" readonly plaintext></b-form-input>
                </b-form-group>
                <b-form-group label="账户类型" horizontal>
                </b-form-group>
                <b-form-group label="昵称" horizontal>
                  <b-form-input v-model="newNickname" type="text"></b-form-input>
                </b-form-group>
                <b-form-group label="创建时间" horizontal>
                  <b-form-input :value="currentUser.created_at | formatUnixEpoch" readonly plaintext></b-form-input>
                </b-form-group>
                <b-button type="submit" class="btn-block" variant="primary">更新</b-button>
              </b-form>
            </b-col>
          </b-row>
        </b-tab>
        <b-tab title="修改密码" class="main-tab">
          <br>I'm the second tab content
        </b-tab>
        <b-tab title="SSH 公钥" class="main-tab">
          <br>I'm the second tab content
        </b-tab>
        <b-tab title="访问令牌" class="main-tab">
          <b-table :items="tokens" :fields="tokenFields">
            <template slot="created_at" slot-scope="data">
              {{data.item.created_at | formatUnixEpoch}}
            </template>
            <template slot="description" slot-scope="data">
              {{data.item.description | formatUserAgent}}
            </template>
            <template slot="action" slot-scope="data">
              <b-link href="#" class="destroy-link" v-if="data.item.id != currentToken.id" @click="onDeleteTokenClick(data.item.id)">删除</b-link>
              <span class="text-muted" v-if="data.item.id == currentToken.id">当前</span>
            </template>
          </b-table>
       </b-tab>
      </b-tabs>
    </b-col>
  </b-row>
</template>

<script>
import {mapGetters, mapState} from 'vuex'
import store from '@/store'

export default {
  name: 'Profile',
  data () {
    return {
      newNickname: null,
      tokenFields: {
        id: {
          label: 'ID',
          sortable: true
        },
        description: {
          label: '浏览器'
        },
        created_at: {
          label: '创建时间'
        },
        action: {
          label: '    '
        }
      }
    }
  },
  computed: {
    ...mapState(['currentToken', 'currentUser', 'tokens']),
    ...mapGetters(['isLoggedIn', 'isLoggedInAsAdmin'])
  },
  mounted () {
    this.newNickname = store.state.currentUser.nickname
    this.$apiListTokens()
  },
  methods: {
    onUpdateSubmit () {
      this.$apiUpdateCurrentUserNickname(this.newNickname).then(() => {
        this.$notify({
          type: 'success',
          title: '昵称修改成功'
        })
      })
    },
    onDeleteTokenClick (id) {
      this.$apiDeleteToken(id).then(() => {
        this.$apiListTokens()
      })
    }
  }
}
</script>

<style scoped>
  #main-row {
    margin-top: 2rem;
  }

  .main-tab {
    padding-top: 1.2rem;
  }
</style>
