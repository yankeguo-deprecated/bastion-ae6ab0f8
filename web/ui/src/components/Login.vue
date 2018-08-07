<template>
  <b-row id="main-row">
    <b-col md="4" offset-md="4">
      <b-card title="登录">
        <b-form @submit="onSubmit">
          <b-form-group label="用户名">
            <b-form-input v-model="form.account" required type="text"></b-form-input>
          </b-form-group>
          <b-form-group label="密码">
            <b-form-input v-model="form.password" required type="password"></b-form-input>
          </b-form-group>
          <b-button type="submit" :disabled="busy" class="btn-block" variant="primary">登录</b-button>
        </b-form>
      </b-card>
    </b-col>
  </b-row>
</template>

<script>
import store from '@/store'

export default {
  name: 'Login',
  data () {
    return {
      form: {
        account: null,
        password: null
      },
      busy: false
    }
  },
  methods: {
    onSubmit () {
      this.busy = true
      this.$http.post('/api/tokens/create', this.form, {emulateJSON: true}).then((res) => {
        store.commit('updateCurrentToken', res.body.token)
        store.commit('updateCurrentUser', res.body.user)
        this.$router.push('/dashboard')
      }, (res) => {
        this.busy = false
        this.$notify({
          type: 'error',
          title: 'API Error',
          text: res.bodyText,
          duration: 2000
        })
      })
    }
  }
}
</script>

<style scoped>
  #main-row {
    margin-top: 2rem;
  }
</style>
