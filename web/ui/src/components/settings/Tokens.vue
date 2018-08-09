<template>
  <b-table :items="tokens" :fields="fields">
    <template slot="created_at" slot-scope="data">
      {{data.item.created_at | formatUnixEpoch}}
    </template>
    <template slot="viewed_at" slot-scope="data">
      {{data.item.viewed_at | formatUnixEpoch}}
    </template>
    <template slot="description" slot-scope="data">
      {{data.item.description | formatUserAgent}}
    </template>
    <template slot="action" slot-scope="data">
      <b-link href="#" class="destroy-link" v-if="data.item.id != currentToken.id"
              @click="onDeleteClick(data.item.id)">删除
      </b-link>
      <span class="text-muted" v-if="data.item.id == currentToken.id">(当前令牌)</span>
    </template>
  </b-table>
</template>

<script>
import {mapGetters, mapState} from 'vuex'

export default {
  name: 'Tokens',
  mounted () {
    this.$apiListTokens()
  },
  data () {
    return {
      fields: [
        {key: 'id', label: 'ID', sortable: true},
        {key: 'description', label: '详情'},
        {key: 'created_at', label: '创建时间', sortable: true},
        {key: 'viewed_at', label: '最后使用时间', sortable: true},
        {key: 'action', label: '操作'}
      ]
    }
  },
  computed: {
    ...mapState(['currentUser', 'currentToken', 'tokens'])
  },
  methods: {
    onDeleteClick (id) {
      this.$apiDeleteToken(id)
    }
  }
}
</script>

<style scoped>

</style>
