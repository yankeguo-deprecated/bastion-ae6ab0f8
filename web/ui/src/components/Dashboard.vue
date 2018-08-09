<template>
<div></div>
</template>

<script>
export default {
  name: 'Dashboard',
  mounted () {
    this.$apiGetCurrentUser()
    this.$apiGetCurrentUserGrantItems()
  }
}
</script>

<style scoped>
</style>
