<template>
  <b-row>
    <b-col md="6" lg="4">
      <b-card>
        <b-form @submit="onSubmit">
          <b-form-group label="名称:" label-class="text-right" description="仅作为备注，默认使用 SSH 公钥备注名" horizontal>
            <b-form-input v-model="form.name" placeholder="输入公钥名称" type="text"></b-form-input>
          </b-form-group>
          <b-form-group label="SSH 公钥:" label-class="text-right" horizontal>
            <b-form-textarea v-model="form.publicKey" :rows="3" placeholder="ssh-rsa AAAAB3N24d..."></b-form-textarea>
          </b-form-group>
          <b-button type="submit" class="btn-block" :disabled="busy" variant="primary">添加</b-button>
        </b-form>
      </b-card>
    </b-col>
    <b-col md="12" lg="12">
      <b-table striped :items="keys" :fields="fields" class="mt-4">
        <template slot="created_at" slot-scope="data">
          {{data.item.created_at | formatUnixEpoch}}
        </template>
        <template slot="viewed_at" slot-scope="data">
          {{data.item.viewed_at | formatUnixEpoch}}
        </template>
        <template slot="fingerprint" slot-scope="data">
          <code>{{data.item.fingerprint}}</code>
        </template>
        <template slot="action" slot-scope="data">
          <b-link href="#" class="destroy-link" v-if="data.item.source != 'sandbox' && data.item.fingerprint != fingerprintToDelete"
                  @click="onDeleteClick(data.item.fingerprint)">删除
          </b-link>
          <b-link href="#" class="destroy-link" v-if="data.item.source != 'sandbox' && data.item.fingerprint == fingerprintToDelete"
                  @click="onDeleteConfirmClick(data.item.fingerprint)">确认删除
          </b-link>
          <span class="text-muted" v-if="data.item.source == 'sandbox'">(沙箱 SSH 公钥)</span>
        </template>
      </b-table>
    </b-col>
  </b-row>
</template>

<script>
/* eslint-disable standard/object-curly-even-spacing,no-unused-vars */

import {mapGetters, mapState} from 'vuex'
export default {
  name: 'Keys',
  data () {
    return {
      fields: [
        { key: 'name', label: '名称', thClass: 'text-center'},
        { key: 'fingerprint', label: '指纹', thClass: 'text-center'},
        { key: 'created_at', label: '添加时间', thClass: 'text-center', tdClass: 'text-center'},
        { key: 'viewed_at', label: '最后使用时间', thClass: 'text-center', tdClass: 'text-center'},
        { key: 'action', label: '    ', thClass: 'text-center', tdClass: 'action-cell'}
      ],
      form: {
        name: '',
        publicKey: ''
      },
      fingerprintToDelete: '',
      busy: false
    }
  },
  computed: {
    ...mapState(['keys'])
  },
  mounted () {
    this.$apiListKeys()
  },
  methods: {
    onSubmit () {
      this.busy = true
      this.$apiCreateKey(this.form).then(() => {
        this.busy = false
        this.form.name = ''
        this.form.publicKey = ''
      }, () => {
        this.busy = false
      })
    },
    onDeleteClick (fingerprint) {
      this.fingerprintToDelete = fingerprint
    },
    onDeleteConfirmClick (fingerprint) {
      this.$apiDestroyKey(fingerprint)
    }
  }
}
</script>

<style scoped></style>
