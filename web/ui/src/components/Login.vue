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
      this.$apiLogin(this.form).then(() => {
        this.busy = false
        this.$router.push('/dashboard')
      }, () => {
        this.busy = false
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
