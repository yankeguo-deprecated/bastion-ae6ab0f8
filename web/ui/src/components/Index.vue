<template>
  <div></div>
</template>

<script>
import {mapGetters} from 'vuex'

export default {
  name: 'Index',
  mounted () {
    if (!this.isLoggedIn) {
      this.$router.push('/login')
    } else {
      this.$router.push('/dashboard')
    }
  },
  computed: {
    ...mapGetters(['isLoggedIn'])
  }
}
</script>

<style scoped></style>
