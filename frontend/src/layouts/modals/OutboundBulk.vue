<template>
  <v-dialog transition="dialog-bottom-transition" width="800" :model-value="visible">
    <v-card class="rounded-lg">
      <v-card-title>
        {{ $t('actions.addbulk') }} {{ $t('objects.outbound') }}
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text style="padding: 0 16px; overflow-y: scroll;">
        <v-row v-if="outbounds.length==0">
          <v-col cols="12">
            <v-text-field v-model="link"
              dir="ltr"
              :label="$t('client.sub')"
              placeholder="http[s]://<domain>[:]<port>/<path>"
              hide-details />
          </v-col>
          <v-col cols="12">
            <v-checkbox v-model="addUrlTest" :label="$t('out.addUrlTest')" />
          </v-col>
          <v-col cols="12" align="center">
            <v-btn hide-details variant="tonal" :loading="loading" @click="linkCheck">{{ $t('submit') }}</v-btn>
          </v-col>
        </v-row>
        <v-data-table
          v-if="outbounds.length>0"
          :items="outbounds"
          :loading="loading"
          :items-per-page="0"
          hide-default-footer
          density="compact"
          :headers="[
            { value: 'check' },
            { title: $t('type'), value: 'type' },
            { title: $t('objects.tag'), value: 'tag' },
            { title: $t('out.addr'), value: 'server' },
            { title: $t('objects.tls'), value: 'tls' }
          ]"
        >
          <template v-slot:[`item.check`]="{ index }">
            <v-icon color="success" icon="mdi-check" v-if="outChecks[index]==1" />
            <v-icon color="error" icon="mdi-close" v-else-if="outChecks[index]==2" />
            <v-progress-circular v-else-if="outChecks[index]==3" indeterminate />
            <v-icon v-else icon="mdi-help"></v-icon>
          </template>
          <template v-slot:[`item.type`]="{ item }">
            {{ item.type }}
          </template>
          <template v-slot:[`item.tag`]="{ item }">
            {{ item.tag }}
          </template>
          <template v-slot:[`item.tls`]="{ item }">
            {{ Object.hasOwn(item,'tls') ? $t(item.tls?.enabled ? 'enable' : 'disable') : '-' }}
          </template>
          <template v-slot:[`item.server`]="{ item }">
            {{ item.server }}{{ item.server_port ? ':' + item.server_port : '' }}
          </template>
        </v-data-table>
      </v-card-text>
      <v-card-actions>
        <v-spacer></v-spacer>
        <v-btn color="primary" variant="outlined" @click="closeModal">{{ $t('actions.close') }}</v-btn>
        <v-btn color="primary" variant="tonal" :loading="loading" :disabled="outbounds.length==0" @click="saveChanges">{{ $t('actions.save') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import HttpUtils from '@/plugins/httputil'
import RandomUtil from '@/plugins/randomUtil';
import Data from '@/store/modules/data'
import { createOutbound, Outbound } from '@/types/outbounds'

export default {
  props: ['visible', 'outboundTags'],
  emits: ['close'],
  data() {
    return {
      loading: false,
      link: "",
      outbounds: <Outbound[]>[],
      outChecks: <number[]>[],
      addUrlTest: false,
    }
  },
  methods: {
    resetData() {
      this.outbounds = []
      this.outChecks = []
      this.link = ""
      this.addUrlTest = false
      this.loading = false
    },
    closeModal() {
      this.resetData()
      this.$emit('close')
    },
    async linkCheck() {
      this.loading = true
      this.outbounds = []
      const msg = await HttpUtils.post('api/subConvert', { link: this.link })
      if (msg.success) {
        if (msg.obj?.length>0) {
          msg.obj.forEach((o:any, index:number) => {
            if (this.newOutboundTags.includes(o.tag)) o.tag = o.tag + "-" + (index+1)
            this.outbounds.push(createOutbound(o.type, o))
            this.outChecks.push(0)
          })
          if (this.addUrlTest) {
            const urlTestTsg = "urltest-" + RandomUtil.randomSeq(3)
            this.outbounds.push(createOutbound("urltest", {
              tag: urlTestTsg,
              outbounds: this.outbounds.map((o:Outbound) => o.tag),
              interrupt_exist_connections: false,
              interval: "30s"
            }))
          }
        }
      }
      this.loading = false
    },
    async saveChanges() {
      if (!this.$props.visible) return
      // check duplicate tag
      this.outbounds.forEach((o:Outbound, index:number) => {
        const isDuplicatedTag = Data().checkTag("outbound",0, o.tag)
        this.outChecks[index] = isDuplicatedTag ? 2 : 0
      })

      // save data
      this.loading = true
      this.outbounds.forEach(async (o:Outbound, index:number) => {
        if (this.outChecks[index] == 2) return
        this.outChecks[index] = 3
        const success = await Data().save("outbounds",  "new", o)
        if (success) this.outChecks[index] = 1
        else this.outChecks[index] = 2
      })
      this.loading = false
    }
  },
  computed: {
    newOutboundTags(): string[] {
      return this.outbounds.map((o:Outbound) => o.tag)
    }
  },
  watch: {
    visible(v) {
      if (v) {
        this.resetData()
      }
    },
  },
}
</script>
