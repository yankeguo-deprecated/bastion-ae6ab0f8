<template>
  <b-row class="mt-4">
    <b-col md="4" lg="3">
      <b-card header="添加/更新服务器" header-tag="b">
        <b-form @submit="onCreateSubmit">
          <b-form-group label="主机名" description="输入符合 Linux 规则的主机名">
            <b-form-input v-model="form.hostname" placeholder="请输入主机名" type="text"></b-form-input>
          </b-form-group>
          <b-form-group label="地址" description="输入服务器的 IP 地址，如果 SSHD 运行在非 22 端口，需要额外注明">
            <b-form-input v-model="form.address" placeholder="请输入地址" type="text"></b-form-input>
          </b-form-group>
          <b-button type="submit" :disabled="busy" variant="success" class="btn-block">
            <i class="fa fa-pencil-square-o" aria-hidden="true"></i> 添加/更新
          </b-button>
        </b-form>
      </b-card>
    </b-col>
    <b-col md="8" lg="9">
      <b-row>
        <b-col>
          <b-card no-body header="服务器列表" header-tag="b">
            <b-card-body>
              <b-form @reset="onReset" inline>
                <b-input v-model="search" class="mb-2 mr-sm-2 mb-sm-0" placeholder="搜索主机名或地址"/>
                <b-button type="reset" :disabled="search == ''" variant="outline-danger">
                  <i class="fa fa-ban" aria-hidden="true"></i> 清除
                </b-button>
              </b-form>
            </b-card-body>
            <b-table striped :items="filteredNodes" :fields="fields" class="mb-0" :show-empty="true" empty-text="无">
              <template slot="created_at" slot-scope="data">
                {{data.item.created_at | formatUnixEpoch}}
              </template>
              <template slot="source" slot-scope="data">
                {{data.item.source}}
              </template>
              <template slot="action" slot-scope="data">
                <b-link href="#" class="text-danger"
                        v-if="data.item.source == 'manual' && hostnameToDelete != data.item.hostname"
                        @click="onDeleteClick(data.item.hostname)"><i class="fa fa-trash" aria-hidden="true"></i> 删除
                </b-link>
                <b-link href="#" class="text-danger"
                        v-if="data.item.source == 'manual' && hostnameToDelete == data.item.hostname"
                        @click="onDeleteConfirmClick(data.item.hostname)"><i class="fa fa-trash" aria-hidden="true"></i> 确认删除
                </b-link>
              </template>
            </b-table>
          </b-card>
        </b-col>
      </b-row>
    </b-col>
  </b-row>
</template>

<script>
/* eslint-disable standard/object-curly-even-spacing */

import {mapState} from 'vuex'

export default {
  name: 'Servers',
  data () {
    return {
      fields: [
        {
          key: 'hostname',
          label: '主机名',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'address',
          label: '地址',
          sortable: true,
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'source',
          label: '来源',
          thClass: 'text-center',
          tdClass: 'text-center'
        },
        {
          key: 'created_at',
          label: '创建日期',
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
      form: {
        hostname: '',
        address: ''
      },
      search: '',
      busy: false,
      hostnameToDelete: ''
    }
  },
  mounted () {
    this.$apiListNodes()
  },
  computed: {
    ...mapState(['nodes']),
    filteredNodes () {
      if (this.search && this.search.length > 0) {
        return this.nodes.filter(n => {
          console.log(n)
          return (
            n.hostname.includes(this.search) || n.address.includes(this.search)
          )
        })
      } else {
        return this.nodes
      }
    }
  },
  methods: {
    onDeleteClick (hostname) {
      this.hostnameToDelete = hostname
    },
    onDeleteConfirmClick (hostname) {
      this.hostnameToDelete = ''
      this.$apiDestroyNode(hostname)
    },
    onReset () {
      this.search = ''
    },
    onCreateSubmit () {
      this.$apiCreateNode(this.form)
    }
  }
}
</script>

<style scoped></style>
