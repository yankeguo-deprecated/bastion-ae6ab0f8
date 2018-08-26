<template>
  <b-row class="mt-3">
    <b-col>
      <b-row>
        <b-col>
          <b-pagination-nav size="md" base-url="#/sessions?page=" :number-of-pages="number_of_session_pages"
                            v-model="currentPage" align="center"></b-pagination-nav>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <b-table striped :items="items" :fields="fields">
            <template slot="account" slot-scope="data">
              <b-link :to="{name: 'UserDetail', params: {account: data.item.account}}">{{data.item.account}}</b-link>
            </template>
            <template slot="command" slot-scope="data">
              <code v-if="data.item.command">{{data.item.command}}</code>
              <code v-if="!data.item.command">(shell)</code>
            </template>
            <template slot="created_at" slot-scope="data">
              {{data.item.created_at | formatUnixEpoch}}
            </template>
            <template slot="finished_at" slot-scope="data">
              {{data.item.finished_at | formatUnixEpoch}}
            </template>
            <template slot="action" slot-scope="data">
              <b-link @click="onReplayClick(data.item.id)" class="text-success" v-if="data.item.is_recorded"><i
                class="fa fa-search" aria-hidden="true"></i> 查看录像
              </b-link>
            </template>
          </b-table>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <b-pagination-nav size="md" base-url="#/sessions?page=" :number-of-pages="number_of_session_pages"
                            v-model="currentPage" align="center"></b-pagination-nav>
        </b-col>
      </b-row>
    </b-col>
  </b-row>
</template>

<script>
import {mapState} from 'vuex'

export default {
  name: 'Sessions',
  data () {
    return {
      currentPage: 1,
      numberOfPages: 9999999999,
      items: [],
      fields: [
        {
          key: 'id',
          label: '编号',
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'account',
          label: '用户',
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'command',
          label: '命令',
          thClass: 'text-center'
        },
        {
          key: 'created_at',
          label: '开始时间',
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'finished_at',
          label: '结束时间',
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'action',
          label: '    ',
          thClass: 'text-center',
          tdClass: 'text-center'
        }
      ]
    }
  },
  computed: {
    ...mapState(['number_of_session_pages'])
  },
  mounted () {
    this.currentPage = Number.parseInt(this.$route.query.page) || 1
    this.listSessions(this.currentPage)
  },
  methods: {
    listSessions (page) {
      this.items = []
      this.$apiListSessions({skip: (page - 1) * 100, limit: 100}).then((res) => {
        this.$store.commit('setNumberOfSessionPages', Math.ceil(res.body.total / 100))
        this.items = res.body.sessions || []
      })
    },
    onReplayClick (id) {
      window.open(`/replays/${id}`, '_blank')
    }
  }
}
</script>

<style scoped>
</style>
