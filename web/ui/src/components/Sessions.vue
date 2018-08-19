<template>
  <b-row class="mt-3">
    <b-col>
      <b-row>
        <b-col>
          <b-pagination size="md" :total-rows="total" v-model="currentPage" :per-page="limit" @change="onPageChanged"
                        align="center"></b-pagination>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <b-table :items="items" :fields="fields">
            <template slot="account" slot-scope="data">
              <b-link :to="{name: 'UserDetail', params: {account: data.item.account}}">{{data.item.account}}</b-link>
            </template>
            <template slot="command" slot-scope="data">
              <code>{{data.item.command}}</code>
            </template>
            <template slot="created_at" slot-scope="data">
              {{data.item.created_at | formatUnixEpoch}}
            </template>
            <template slot="finished_at" slot-scope="data">
              {{data.item.finished_at | formatUnixEpoch}}
            </template>
            <template slot="action" slot-scope="data">
              <b-link href="#" class="text-success" v-if="data.item.is_recorded"
                      @click="onViewReplayClicked(data.item.id)"><i class="fa fa-search" aria-hidden="true"></i> 查看录像
              </b-link>
            </template>
          </b-table>
        </b-col>
      </b-row>
      <b-row>
        <b-col>
          <b-pagination size="md" :total-rows="total" v-model="currentPage" :per-page="limit" @change="onPageChanged"
                        align="center"></b-pagination>
        </b-col>
      </b-row>
    </b-col>
  </b-row>
</template>

<script>
export default {
  name: 'Sessions',
  data () {
    return {
      limit: 100,
      total: 0,
      currentPage: 1,
      items: [],
      fields: [
        {
          key: 'id',
          label: '编号'
        },
        {
          key: 'account',
          label: '用户'
        },
        {
          key: 'command',
          label: '命令'
        },
        {
          key: 'created_at',
          label: '开始时间'
        },
        {
          key: 'finished_at',
          label: '结束时间'
        },
        {
          key: 'action',
          label: '    '
        }
      ]
    }
  },
  mounted () {
    this.listSessions(this.currentPage)
  },
  methods: {
    listSessions (page) {
      this.items = []
      this.$apiListSessions({skip: (page - 1) * this.limit, limit: this.limit}).then((res) => {
        this.total = res.body.total
        this.items = res.body.sessions || []
      })
    },
    onPageChanged (page) {
      this.listSessions(page)
    },
    onViewReplayClicked (id) {
      alert(id)
    }
  }
}
</script>

<style scoped>
</style>
