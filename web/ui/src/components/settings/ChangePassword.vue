<template>
  <b-row>
    <b-col md="6" lg="4">
      <b-card>
        <b-form @submit="onSubmit">
          <b-form-group label="旧密码" horizontal>
            <b-form-input v-model="form.oldPassword" type="password"></b-form-input>
          </b-form-group>
          <b-form-group label="新密码" horizontal>
            <b-form-input v-model="form.newPassword" type="password"></b-form-input>
          </b-form-group>
          <b-form-group label="重复密码" horizontal>
            <b-form-input v-model="form.repPassword" type="password"></b-form-input>
          </b-form-group>
          <b-button type="submit" class="btn-block" variant="primary">修改密码</b-button>
        </b-form>
      </b-card>
    </b-col>
  </b-row>
</template>

<script>
export default {
  name: 'ChangePassword',
  data () {
    return {
      form: {
        oldPassword: '',
        newPassword: '',
        repPassword: ''
      }
    }
  },
  methods: {
    onSubmit () {
      if (this.form.oldPassword.length === 0 ||
      this.form.newPassword.length === 0) {
        this.$notify({
          type: 'warn',
          title: '输入错误',
          text: '旧密码和新密码不能为空'
        })
        return
      }
      if (this.form.newPassword !== this.form.repPassword) {
        this.$notify({
          type: 'warn',
          title: '输入错误',
          text: '重复密码和新密码不匹配'
        })
        return
      }
      this.$apiUpdatePassword(this.form).then(() => {
        this.form.oldPassword = ''
        this.form.newPassword = ''
        this.form.repPassword = ''
      })
    }
  }
}
</script>

<style scoped>
</style>
