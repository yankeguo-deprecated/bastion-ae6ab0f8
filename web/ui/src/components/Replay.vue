<template>
  <b-row>
    <div class="info-bar">
      <b-nav>
        <b-nav-item disabled>{{user.account}}({{user.nickname}})</b-nav-item>
        <b-nav-item disabled>{{session.command}}</b-nav-item>
        <b-nav-item disabled>{{session.created_at | formatUnixEpoch}} => {{session.finished_at | formatUnixEpoch}}</b-nav-item>
      </b-nav>
    </div>
    <div id="terminal"></div>
    <div class="control-bar"></div>
  </b-row>
</template>

<script>
export default {
  name: 'Replay',
  data () {
    return {
      session: {},
      user: {}
    }
  },
  mounted () {
    this.$apiGetSession(this.$route.params.id).then((res) => {
      this.session = res.body.session
      this.user = res.body.user
    })
  }
}
</script>

<style scoped>
  div.info-bar {
    position: fixed;
    top: 0;
    height: 56px;
    width: 100%;
    left: 0;
    background-color: #ced6e0;
    padding: 8px;
  }

  div#terminal {
    position: fixed;
    left: 0;
    width: 100%;
    top: 56px;
    bottom: 56px;
    background-color: black;
  }

  div.control-bar {
    position: fixed;
    bottom: 0;
    height: 56px;
    width: 100%;
    left: 0;
    padding: 8px;
    background-color: #ced6e0;
  }

</style>
