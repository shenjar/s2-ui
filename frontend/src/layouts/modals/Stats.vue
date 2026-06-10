<template>
  <v-dialog transition="dialog-bottom-transition" width="800">
    <v-card class="rounded-lg" :loading="loading">
      <v-card-title>
        <v-row>
          <v-col cols="auto">
            {{ $t('stats.graphTitle') }}
          </v-col>
          <v-spacer></v-spacer>
          <v-col cols="auto"><v-icon icon="mdi-close" @click="$emit('close')"></v-icon></v-col>
        </v-row>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text style="padding: 0 16px;">
        <div style="text-align: center; margin: 5px;">
          {{ $t('objects.' + resource) + " : " + tag }}
        </div>
        <v-radio-group v-model="limit" @update:model-value="onLimitChange" density="compact" :loading="loading" inline hide-details>
          <v-radio v-for="p in periods" :label="p.title" :value="p.value"></v-radio>
        </v-radio-group>
          <v-container id="container" style="height:40vh;">
            <v-skeleton-loader
            class="mx-auto border"
            width="95%"
            type="image"
            v-if="loading"
          ></v-skeleton-loader>
          <template v-else>
            <v-alert :text="$t('noData')" type="warning" variant="outlined" v-if="alert"></v-alert>
            <Line v-if="loaded" :data="usage" :options="<any>options" />
          </template>
        </v-container>
      </v-card-text>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { i18n } from '@/locales'
import HttpUtils from '@/plugins/httputil'
import { HumanReadable } from '@/plugins/utils'
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js'
import { ref } from 'vue'
import { Line } from 'vue-chartjs'
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  Filler
)
ChartJS.defaults.font.family = 'Vazirmatn'
export default {
  components: {
    Line
  },
  props: ['visible','resource','tag'],
  data() {
    return {
      loading: false,
      loaded: false,
      alert: false,
      intervalId: <any>0,
      limit: 'hour',
      periods: [
        { value: 'hour',  title: i18n.global.n(1) + i18n.global.t('date.h') },
        { value: 'day',   title: i18n.global.n(1) + i18n.global.t('date.d') },
        { value: 'month', title: i18n.global.n(30) + i18n.global.t('date.d') },
      ],
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: {
          intersect: false,
          mode: 'index',
        },
        elements: {
          point: { pointStyle: 'crossRot' }
        },
        plugins: {
          tooltip: {
            callbacks: {
              label: (ctx:any) => {
                return ctx.dataset.label + ': ' + HumanReadable.sizeFormat(Number(ctx.raw) || 0)
              },
              footer: (items:any[]) => {
                return HumanReadable.sizeFormat(items.reduce((acc, c) => acc + (Number(c.raw) || 0), 0))
              }
            }
          }
        },
        scales: {
          x: {
            ticks: { maxTicksLimit: 10, maxRotation: 0, autoSkip: true }
          },
          y: {
            grid: {
              color: '#777777',
            },
            beginAtZero: true,
            ticks: {
              callback: function(label:any, index: number) {
                return label == 0 ? 0 : HumanReadable.sizeFormat(label,0)
              },
              count: 10
            }
          }
        }
      },
      usage: ref(<any>{}),
    }
  },
  methods: {
    onLimitChange(value: string | null) {
      if (value == null) return
      this.limit = value
      this.loadData()
    },
    async loadData() {
      this.loading = true
      const data = await HttpUtils.get('api/stats', { resource: this.resource, tag: this.tag, period: this.limit })
      if (data.success && data.obj) {
        const obj = <any[]>data.obj
        if (obj.length == 0) {
          this.usage = { labels: [], datasets: [] }
          this.loaded = false
          this.alert = true
          this.loading = false
          return
        }
        const l = String(i18n.global.locale) == 'fa' ? "fa-IR" : "en-US"
        // The backend returns exactly one up row and one down row per time
        // bucket (empty buckets included with traffic=0), already sorted by
        // time. We split them into two series keyed by timestamp.
        const upByTime: Record<number, number> = {}
        const downByTime: Record<number, number> = {}
        const timeSet = new Set<number>()
        for (const o of obj) {
          const t = Number(o.dateTime)
          if (!t) continue
          timeSet.add(t)
          if (o.direction) {
            upByTime[t] = Number(o.traffic) || 0
          } else {
            downByTime[t] = Number(o.traffic) || 0
          }
        }
        const times = Array.from(timeSet).sort((a, b) => a - b)
        if (times.length == 0) {
          this.usage = { labels: [], datasets: [] }
          this.loaded = false
          this.alert = true
          this.loading = false
          return
        }
        const labels = times.map(t => this.genLabel(t * 1000, l))
        const uplinkData = times.map(t => upByTime[t] ?? 0)
        const downlinkData = times.map(t => downByTime[t] ?? 0)
        this.usage = {
          labels: labels,
          datasets: [
            {
              label: i18n.global.t('stats.upload'),
              backgroundColor: 'rgba(255, 165, 0, 0.4)',
              borderColor: 'rgba(255, 165, 0)',
              fill: true,
              data: uplinkData
            },
            {
              label: i18n.global.t('stats.download'),
              backgroundColor: 'rgba(0, 128, 0, 0.2)',
              borderColor: 'rgba(0, 128, 0)',
              fill: true,
              data: downlinkData
            }
          ],
        }
        this.loaded = true
        this.alert = false
      } else {
        this.alert = true
        this.loaded = false
      }
      this.loading = false
    },
    genLabel(step: number, locale: string) {
      const d = new Date(step)
      if (this.limit === 'month') {
        return d.toLocaleDateString(locale, { month: '2-digit', day: '2-digit' })
      }
      return d.toLocaleTimeString(locale, { hour: '2-digit', minute: '2-digit', hour12: false })
    },
  },
  watch: {
    visible(v) {
      if (v) {
        this.limit = 'hour'
        this.loadData()
        this.intervalId = setInterval(() => {
          this.loadData()
        }, 10000)
      } else {
        this.loaded = false
        this.alert = false
        this.usage.labels = []
        if (this.usage.datasets) {
          this.usage.datasets[0].data = []
          this.usage.datasets[1].data = []
        }
        if (this.intervalId && this.intervalId != 0) {
          clearInterval(this.intervalId)
        }
      }
    }
  }
}
</script>
