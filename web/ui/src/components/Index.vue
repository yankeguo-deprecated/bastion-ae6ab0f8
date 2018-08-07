<template>
  <div></div>
</template>

<script>

import store from '@/store'

export default {
  name: 'Index',
  mounted () {
    if (!store.getters.isLoggedIn) {
      this.$router.push('/login')
    } else {
      this.$router.push('/dashboard')
    }
  }
}
</script>

<style scoped>
</style>
