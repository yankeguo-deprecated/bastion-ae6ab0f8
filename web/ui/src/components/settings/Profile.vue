<template>
  <b-row>
    <b-col md="6" lg="4">
      <b-card>
        <b-form @submit="onUpdateSubmit">
          <b-form-group label="账户名" horizontal>
            <b-form-input :value="currentUser.account" readonly plaintext></b-form-input>
          </b-form-group>
          <b-form-group label="账户类型" horizontal>
            <b-form-input :value="userType" readonly plaintext></b-form-input>
          </b-form-group>
          <b-form-group label="昵称" horizontal>
            <b-form-input v-model="nickname" type="text"></b-form-input>
          </b-form-group>
          <b-form-group label="创建时间" horizontal>
            <b-form-input :value="currentUser.created_at | formatUnixEpoch" readonly plaintext></b-form-input>
          </b-form-group>
          <b-button type="submit" class="btn-block" variant="primary">更新</b-button>
        </b-form>
      </b-card>
    </b-col>
  </b-row>
</template>

<script>
import {mapGetters, mapState} from 'vuex'

export default {
  name: 'Profile',
  data () {
    return {
      nickname: null
    }
  },
  computed: {
    ...mapState(['currentUser']),
    ...mapGetters(['userType'])
  },
  mounted () {
    this.nickname = this.currentUser.nickname
    this.$apiGetCurrentUser()
  },
  methods: {
    onUpdateSubmit () {
      this.$apiUpdateCurrentUserNickname(this.nickname)
    }
  }
}
</script>

<style scoped></style>
