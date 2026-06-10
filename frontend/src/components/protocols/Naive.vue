<template>
  <v-card>
    <v-card-subtitle v-if="direction != 'out_json'">Naive</v-card-subtitle>
    <!-- Inbound -->
    <template v-if="direction === 'in'">
      <v-row>
        <v-col cols="12" sm="6" md="4">
          <Network :data="data" />
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-select
            hide-details
            :label="$t('types.naive.quicCongestion')"
            :items="inbCngs"
            v-model="data.quic_congestion_control"
            @click:clear="delete data.quic_congestion_control"
            clearable>
          </v-select>
        </v-col>
      </v-row>
    </template>
    <!-- Outbound -->
    <template v-if="['out', 'out_json'].includes(direction)">
      <v-row v-if="direction === 'out'">
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            :label="$t('types.un')"
            hide-details
            v-model="data.username">
          </v-text-field>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            :label="$t('types.pw')"
            hide-details
            type="password"
            v-model="data.password">
          </v-text-field>
        </v-col>
      </v-row>
      <v-row>
        <v-col cols="12" sm="6" md="4">
          <v-text-field
            :label="$t('types.naive.insecureConcurrency')"
            type="number"
            min="0"
            hide-details
            v-model.number="insecure_concurrency">
          </v-text-field>
        </v-col>
        <v-col cols="12" sm="6" md="4">
          <v-switch v-model="udpOverTcp" color="primary" :label="$t('types.naive.udpOverTcp')" hide-details></v-switch>
        </v-col>
      </v-row>
      <v-row v-if="direction === 'out'">
        <v-col cols="12" sm="6" md="4">
          <v-switch v-model="data.quic" color="primary" :label="$t('types.naive.quic')" hide-details></v-switch>
        </v-col>
        <v-col cols="12" sm="6" md="4" v-if="data.quic">
          <v-select
            hide-details
            :label="$t('types.naive.quicCongestion')"
            :items="outCngs"
            @click:clear="delete data.quic_congestion_control"
            clearable
            v-model="data.quic_congestion_control">
          </v-select>
        </v-col>
      </v-row>
      <Headers :data="extra_headers" />
    </template>
  </v-card>
</template>

<script lang="ts">
import Network from '@/components/Network.vue'
import Headers from '@/components/Headers.vue'

export default {
  props: ['data', 'direction'],
  data() {
    return {
      inbCngs: [
        { title: 'BBR', value: 'bbr'},
        { title: 'BBR Standard', value: 'bbr_standard'},
        { title: 'BBRv2', value: 'bbr2'},
        { title: 'BBRv2 variant', value: 'bbr2_variant'},
        { title: 'Cubic', value: 'cubic'},
        { title: 'New Reno', value: 'reno'},
      ],
      outCngs: [
        { title: 'BBR', value: 'bbr'},
        { title: 'BBR2', value: 'bbr2'},
        { title: 'Cubic', value: 'cubic'},
        { title: 'Reno', value: 'reno'},
      ],
    }
  },
  computed: {
    udpOverTcp: {
      get(): boolean {
        const d = this.$props.data
        return d?.udp_over_tcp === true || (d?.udp_over_tcp && (d.udp_over_tcp as any)?.enabled)
      },
      set(v: boolean) {
        if (v) {
          this.$props.data.udp_over_tcp = { enabled: true }
        } else {
          this.$props.data.udp_over_tcp = false
        }
      }
    },
    insecure_concurrency: {
      get(): number { return this.$props.data?.insecure_concurrency ?? 0 },
      set(v: number) {
        this.$props.data.insecure_concurrency = v > 0 ? v : undefined
      }
    },
    extra_headers(): any {
      const d = this.$props.data
      return new Proxy({}, {
        get(_, prop) {
          if (prop === 'headers') return d?.extra_headers ?? {}
          return undefined
        },
        set(_, prop, value) {
          if (prop === 'headers') {
            d.extra_headers = value
            return true
          }
          return false
        }
      })
    },
  },
  components: { Network, Headers }
}
</script>
