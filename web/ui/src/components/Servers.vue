<template>
  <b-row class="mt-4">
    <b-modal id="modal1" title="添加/更新服务器" :hide-footer="true">
      <b-form @submit="onCreateSubmit">
        <b-form-group label="主机名:" label-class="text-right" description="输入符合 Linux 规则的主机名" horizontal>
          <b-form-input v-model="form.hostname" placeholder="请输入主机名" type="text"></b-form-input>
        </b-form-group>
        <b-form-group label="地址:" label-class="text-right" description="输入服务器的 IP 地址，如果 SSHD 运行在非 22 端口，需要额外注明" horizontal>
          <b-form-input v-model="form.address" placeholder="请输入地址" type="text"></b-form-input>
        </b-form-group>
        <div class="text-right">
          <b-button type="submit" :disabled="busy" variant="success"><i class="fa fa-pencil-square-o" aria-hidden="true"></i> 添加/更新</b-button>
        </div>
      </b-form>
    </b-modal>
    <b-col>
      <b-row>
        <b-col md="4" lg="3">
          <b-form @reset="onReset" inline>
            <b-input v-model="search" class="mb-2 mr-sm-2 mb-sm-0" placeholder="搜索主机名或地址"/>
            <b-button type="reset" :disabled="search == ''" variant="outline-danger"><i class="fa fa-ban" aria-hidden="true"></i> 清除</b-button>
          </b-form>
        </b-col>
        <b-col md="8" lg="9" class="text-right">
          <b-btn variant="success" v-b-modal.modal1><i class="fa fa-pencil-square-o" aria-hidden="true"></i> 添加/更新</b-btn>
        </b-col>
     </b-row>
      <b-row class="mt-4">
        <b-col>
          <b-table striped :items="filteredNodes" :fields="fields">
            <template slot="created_at" slot-scope="data">
              {{data.item.created_at | formatUnixEpoch}}
            </template>
            <template slot="source" slot-scope="data">
              {{data.item.source}}
            </template>
            <template slot="action" slot-scope="data">
              <b-link href="#" class="text-danger"
                      v-if="data.item.source == 'manual' && hostnameToDelete != data.item.hostname"
                      @click="onDeleteClick(data.item.hostname)">删除
              </b-link>
              <b-link href="#" class="text-danger"
                      v-if="data.item.source == 'manual' && hostnameToDelete == data.item.hostname"
                      @click="onDeleteConfirmClick(data.item.hostname)">确认删除
              </b-link>
            </template>
          </b-table>
        </b-col>
      </b-row>
    </b-col>
  </b-row>
</template>

<script>
/* eslint-disable standard/object-curly-even-spacing */

import { mapState } from 'vuex'

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

<style scoped>
</style>
