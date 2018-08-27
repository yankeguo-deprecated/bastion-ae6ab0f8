<template>
  <b-row class="mt-3 mb-3">
    <b-col md="4" lg="2">
      <b-list-group>
        <b-list-group-item to="/settings/profile"><i class="fa fa-user-circle-o" aria-hidden="true"></i> 个人信息</b-list-group-item>
        <b-list-group-item to="/settings/change-password"><i class="fa fa-lock" aria-hidden="true"></i> 修改密码</b-list-group-item>
        <b-list-group-item to="/settings/keys"><i class="fa fa-key" aria-hidden="true"></i> SSH 公钥</b-list-group-item>
        <b-list-group-item to="/settings/tokens"><i class="fa fa-id-card" aria-hidden="true"></i> 访问令牌</b-list-group-item>
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

<style scoped></style>
