<template>
  <b-row>
    <b-col md="6" lg="4">
      <b-card>
        <b-form>
          <b-form-group label="账户名:" label-class="text-right" horizontal>
            <b-form-input :value="currentUser.account" readonly plaintext></b-form-input>
          </b-form-group>
          <b-form-group label="账户类型:" label-class="text-right" horizontal>
            <b-form-input :value="currentUser | formatUserStatus" readonly plaintext></b-form-input>
          </b-form-group>
          <b-form-group label="昵称:" label-class="text-right" horizontal>
            <b-form-input :value="nickname"  readonly plaintext></b-form-input>
          </b-form-group>
          <b-form-group label="创建时间:" label-class="text-right" horizontal>
            <b-form-input :value="currentUser.created_at | formatUnixEpoch" readonly plaintext></b-form-input>
          </b-form-group>
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
  }
}
</script>

<style scoped>
</style>
