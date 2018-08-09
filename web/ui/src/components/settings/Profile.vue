<template>
  <b-row>
    <b-col md="6" lg="4">
      <b-card>
        <b-form @submit="onUpdateSubmit">
          <b-form-group label="账户名:" label-class="text-right" horizontal>
            <b-form-input :value="currentUser.account" readonly plaintext></b-form-input>
          </b-form-group>
          <b-form-group label="账户类型:" label-class="text-right" horizontal>
            <b-form-input :value="currentUser | formatUserStatus" readonly plaintext></b-form-input>
          </b-form-group>
          <b-form-group label="昵称:" label-class="text-right" description="昵称不能大于5个中文字符，不能为空" horizontal>
            <b-form-input v-model="nickname" placeholder="请输入昵称" type="text"></b-form-input>
          </b-form-group>
          <b-form-group label="创建时间:" label-class="text-right" horizontal>
            <b-form-input :value="currentUser.created_at | formatUnixEpoch" readonly plaintext></b-form-input>
          </b-form-group>
          <b-button type="submit" :disabled="busy" class="btn-block" variant="primary">更新</b-button>
        </b-form>
      </b-card>
    </b-col>
  </b-row>
</template>

<script>
import {mapState} from 'vuex'

export default {
  name: 'Profile',
  data () {
    return {
      nickname: null,
      busy: false
    }
  },
  computed: {
    ...mapState(['currentUser'])
  },
  mounted () {
    this.nickname = this.currentUser.nickname
    this.$apiGetCurrentUser()
  },
  methods: {
    onUpdateSubmit () {
      this.busy = true
      this.$apiUpdateCurrentUserNickname(this.nickname).then(() => {
        this.busy = false
      }, () => {
        this.busy = false
      })
    }
  }
}
</script>

<style scoped></style>
