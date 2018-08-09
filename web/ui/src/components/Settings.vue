<template>
  <b-row id="main-row">
    <b-col md="4" lg="2">
      <b-list-group>
        <b-list-group-item to="/settings/profile">个人信息</b-list-group-item>
        <b-list-group-item to="/settings/change-password">修改密码</b-list-group-item>
        <b-list-group-item to="/settings/keys">SSH 公钥</b-list-group-item>
        <b-list-group-item to="/settings/tokens">访问令牌</b-list-group-item>
      </b-list-group>
    </b-col>
    <b-col md="8" lg="10">
      <router-view/>
    </b-col>
  </b-row>
</template>

<script>
export default {
  name: 'Settings'
}
</script>

<style scoped>
#main-row {
  margin-top: 2rem;
}
</style>
