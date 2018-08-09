<template>
  <b-row>
    <b-col md="6" lg="4">
      <b-card>
        <b-form @submit="onSubmit">
          <b-form-group label="名称" horizontal>
            <b-form-input v-model="form.name" type="text"></b-form-input>
          </b-form-group>
          <b-form-group label="SSH 公钥" horizontal>
            <b-form-textarea v-model="form.publicKey" type="text"></b-form-textarea>
          </b-form-group>
          <b-button type="submit" class="btn-block" variant="primary">添加</b-button>
        </b-form>
      </b-card>
    </b-col>
    <b-col md="12" lg="12">
      <b-table :items="keys" :fields="fields" class="main-table">
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
          <b-link href="#" class="destroy-link" v-if="data.item.source != 'sandbox'"
                  @click="onDeleteClick(data.item.fingerprint)">删除
          </b-link>
          <span class="text-muted" v-if="data.item.source == 'sandbox'">(沙箱 SSH 公钥)</span>
        </template>
      </b-table>
    </b-col>
  </b-row>
</template>

<script>
import {mapGetters, mapState} from 'vuex'
export default {
  name: 'Keys',
  data () {
    return {
      fields: [
        { key: 'name', label: '名称'},
        { key: 'fingerprint', label: '指纹'},
        { key: 'created_at', label: '添加时间'},
        { key: 'viewed_at', label: '最后使用时间'},
        { key: 'action', label: '    '}
      ],
      form: {
        name: '',
        publicKey: ''
      }
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
      this.$apiCreateKey(this.form).then(() => {
        this.form.name = ''
        this.form.publicKey = ''
      })
    },
    onDeleteClick (fingerprint) {
      this.$apiDestroyKey(fingerprint)
    }
  }
}
</script>

<style scoped>
  table.main-table {
    margin-top: 2rem;
  }
</style>
