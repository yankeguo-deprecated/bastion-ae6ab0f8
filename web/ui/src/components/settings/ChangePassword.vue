<template>
  <b-row>
    <b-col md="6" lg="4">
      <b-card>
        <b-form @submit="onSubmit">
          <b-form-group label="旧密码:" label-class="text-right" description="输入旧密码" horizontal>
            <b-form-input v-model="form.oldPassword" placeholder="输入旧密码" type="password"></b-form-input>
          </b-form-group>
          <b-form-group label="新密码:" label-class="text-right" description="长度不能小于6" horizontal>
            <b-form-input v-model="form.newPassword" placeholder="输入新密码" type="password"></b-form-input>
          </b-form-group>
          <b-form-group label="重复密码:" label-class="text-right" horizontal>
            <b-form-input v-model="form.repPassword" placeholder="重复输入新密码" type="password"></b-form-input>
          </b-form-group>
          <b-button type="submit" class="btn-block" :disabled="busy" variant="primary"><i class="fa fa-upload" aria-hidden="true"></i> 修改密码</b-button>
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
      },
      busy: false
    }
  },
  methods: {
    onSubmit () {
      if (
        this.form.oldPassword.length === 0 ||
        this.form.newPassword.length === 0
      ) {
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
      this.busy = true
      this.$apiUpdatePassword(this.form).then(
        () => {
          this.busy = false
          this.form.oldPassword = ''
          this.form.newPassword = ''
          this.form.repPassword = ''
        },
        () => {
          this.busy = false
        }
      )
    }
  }
}
</script>

<style scoped>
</style>
