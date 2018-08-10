<template>
  <b-table striped :items="tokens" :fields="fields">
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
      <b-link href="#" class="text-danger" v-if="data.item.id != currentToken.id && data.item.id != tokenToDelete"
              @click="onDeleteClick(data.item.id)"><i class="fa fa-trash" aria-hidden="true"></i> 删除
      </b-link>
      <b-link href="#" class="text-danger" v-if="data.item.id != currentToken.id && data.item.id == tokenToDelete"
              @click="onDeleteConfirmClick(data.item.id)"><i class="fa fa-trash" aria-hidden="true"></i> 确认删除
      </b-link>
      <span class="text-muted" v-if="data.item.id == currentToken.id">(当前)</span>
    </template>
  </b-table>
</template>

<script>
import { mapState } from 'vuex'

export default {
  name: 'Tokens',
  mounted () {
    this.$apiListTokens()
  },
  data () {
    return {
      fields: [
        { key: 'description', label: '详情', thClass: 'text-center' },
        {
          key: 'created_at',
          label: '创建时间',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'viewed_at',
          label: '最后使用时间',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'action',
          label: '    ',
          thClass: 'text-center',
          tdClass: 'action-cell'
        }
      ],
      tokenToDelete: 0
    }
  },
  computed: {
    ...mapState(['currentUser', 'currentToken', 'tokens'])
  },
  methods: {
    onDeleteClick (id) {
      this.tokenToDelete = id
    },
    onDeleteConfirmClick (id) {
      this.tokenToDelete = 0
      this.$apiDeleteToken(id)
    }
  }
}
</script>

<style scoped>
</style>
