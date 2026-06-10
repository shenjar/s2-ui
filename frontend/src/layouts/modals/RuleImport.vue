<template>
  <v-dialog transition="dialog-top-transition" width="800">
    <v-card class="rounded-lg">
      <v-card-title>
        <v-row>
          <v-col>
            {{ $t('rule.import.rulesTitle') }}
          </v-col>
          <v-col cols="auto">
            <v-chip v-if="parsed" size="small" color="primary" variant="tonal">
              {{ parsed.rules?.length ?? 0 }} {{ $t('pages.rules') }} · {{ parsed.rule_set?.length ?? 0 }} {{ $t('rule.ruleset') }}
            </v-chip>
          </v-col>
        </v-row>
      </v-card-title>
      <v-divider></v-divider>
      <v-card-text style="padding: 0 16px; overflow-y: scroll;">
        <v-tabs v-model="tab" @update:modelValue="tabChanged">
          <v-tab value="json">JSON</v-tab>
          <v-tab value="file">{{ $t('rule.import.uploadFile') }}</v-tab>
          <v-tab value="url">{{ $t('rule.import.fromUrl') }}</v-tab>
        </v-tabs>

        <v-window v-model="tab">
          <v-window-item value="json">
            <v-alert variant="text" type="info">{{ $t('rule.import.jsonHint') }}</v-alert>
            <v-textarea
              v-model="rawJson"
              label="JSON"
              variant="outlined"
              rows="12"
              hide-details
              spellcheck="false"
              class="mb-4"
              :error="!!error"></v-textarea>
          </v-window-item>

          <v-window-item value="file">
            <v-alert variant="text" type="info">{{ $t('rule.import.fileJsonHint') }}</v-alert>
            <v-file-input
              :label="$t('rule.import.selectJson')"
              accept=".json"
              variant="outlined"
              hide-details
              prepend-icon="mdi-file-code"
              clearable
              @click:clear="tabChanged"
              @update:modelValue="onFileUpload($event)" />
          </v-window-item>

          <v-window-item value="url">
            <v-alert variant="text" type="info">{{ $t('rule.import.urlHint') }}</v-alert>
            <v-text-field
              v-model="fetchUrl"
              label="URL"
              variant="outlined"
              hide-details
              spellcheck="false"
              placeholder="https://raw.githubusercontent.com/.../rules.json"
              append-icon="mdi-download"
              @keydown.enter="fetchFromUrl"
              @click:append="fetchFromUrl" />
          </v-window-item>
        </v-window>
        <v-alert v-if="error" type="error" variant="text">{{ error }}</v-alert>

        <template v-if="parsed">
          <v-divider class="my-4" />

          <v-alert v-if="hasConflicts" type="warning" variant="tonal" :title="$t('rule.import.conflict')" class="mb-4">
            {{ $t('rule.import.conflictDesc', { rules: existingRulesCount, rulesets: existingRulesetsCount }) }}
            <v-radio-group v-model="mode" hide-details class="mt-2">
              <v-radio value="merge" :label="$t('rule.import.merge')" />
              <v-radio value="replace" :label="$t('rule.import.replace')" />
            </v-radio-group>
          </v-alert>

          <v-alert v-if="parsed.final" type="info" variant="tonal" class="mb-4">
            {{ $t('rule.import.finalOutbound') }}:
            <v-chip size="small" color="secondary" variant="tonal">{{ parsed.final }}</v-chip>
            <v-checkbox v-model="applyFinal" :label="$t('rule.import.applyFinal')" hide-details density="compact" />
          </v-alert>

          <span class="v-card-subtitle">
            {{ $t('pages.rules') }}
            <v-badge v-if="parsed.rules?.length > 0" color="success" :content="parsed.rules?.length" inline />
          </span>
          <v-table v-if="parsed.rules?.length" density="compact" class="mb-4" striped="even">
            <thead>
              <tr><th>#</th><th>{{ $t('type') }}</th><th>{{ $t('admin.action') }}</th><th>{{ $t('objects.outbound') }}</th></tr>
            </thead>
            <tbody>
              <tr v-for="(r, i) in parsed.rules" :key="i">
                <td>{{ (i as number) + 1 }}</td>
                <td>{{ r.type ?? 'simple' }}</td>
                <td>{{ r.action }}</td>
                <td>{{ r.outbound ?? '-' }}</td>
              </tr>
            </tbody>
          </v-table>

          <span class="v-card-subtitle">
            {{ $t('rule.ruleset') }}
            <v-badge v-if="parsed.rule_set?.length > 0" color="success" :content="parsed.rule_set?.length" inline />
            <span v-if="skippedRulesets > 0">
              <v-badge color="warning" :content="skippedRulesets" inline v-tooltip:top="$t('rule.import.skipped')" />
            </span>
          </span>
          <v-table v-if="parsed.rule_set?.length" density="compact" striped="even">
            <thead>
              <tr><th>{{ $t('objects.tag') }}</th><th>{{ $t('ruleset.format') }}</th><th>{{ $t('type') }}</th><th>{{ $t('ruleset.interval') }}</th></tr>
            </thead>
            <tbody>
              <tr v-for="(rs, i) in parsed.rule_set" :key="i"
                :style="mode === 'merge' && existingRulesetTags.includes(rs.tag) ? 'opacity:0.4' : ''">
                <td style="font-family: monospace;">{{ rs.tag }}</td>
                <td>{{ rs.format }}</td>
                <td>{{ rs.type }}</td>
                <td>{{ rs.update_interval ?? '-' }}</td>
              </tr>
            </tbody>
          </v-table>
        </template>
      </v-card-text>
      <v-card-actions>
        <v-btn
          v-if="tab === 'json'"
          @click="parseJson"
          variant="tonal"
          color="success"
          :disabled="rawJson.trim().length === 0"
          >
          {{ $t('rule.import.parse') }}
          <v-icon icon="mdi-magnify" />
        </v-btn>
        <v-spacer />
        <v-btn @click="close" variant="text">{{ $t('actions.close') }}</v-btn>
        <v-btn @click="save" color="primary" variant="flat" :disabled="!parsed">
          {{ $t('actions.save') }}
        </v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
export default {
  props: ["visible", "existingRulesCount", "existingRulesetsCount", "existingRulesetTags"],
  emits: ['save', 'close'],
  data() {
    return {
      tab: 'json',
      rawJson: '',
      fetchUrl: '',
      fetching: false,
      error: '',
      parsed: null as any,
      mode: 'merge' as 'merge' | 'replace',
      applyFinal: false,
    }
  },
  computed: {
    hasConflicts(): boolean {
      return this.existingRulesCount > 0 || this.existingRulesetsCount > 0
    },
    skippedRulesets(): number {
      if (!this.parsed?.rule_set) return 0
      const existing = new Set(this.existingRulesetTags)
      return this.parsed.rule_set.filter((rs: any) => existing.has(rs.tag)).length
    },
  },
  methods: {
    tabChanged() {
      this.rawJson = ''
      this.fetchUrl = ''
      this.error = ''
      this.parsed = null
      this.mode = this.hasConflicts ? 'merge' : 'replace'
      this.applyFinal = false
    },
    extractRouteBlock(obj: any): any {
      if (obj?.route && (obj.route.rules || obj.route.rule_set)) return obj.route
      if (obj?.rules || obj?.rule_set) return obj
      return null
    },
    setParsed(block: any) {
      this.parsed = block
      this.mode = this.hasConflicts ? 'merge' : 'replace'
      this.applyFinal = false
    },
    reset() {
      this.tab = 'json'
      this.tabChanged()
    },
    parseJson() {
      this.error = ''
      this.parsed = null
      try {
        const block = this.extractRouteBlock(JSON.parse(this.rawJson))
        if (!block) {
          this.error = this.$t('rule.import.errNoArrays')
          return
        }
        this.setParsed(block)
      } catch (e: any) {
        this.error = this.$t('rule.import.errJsonParse', { message: e.message })
      }
    },
    async fetchFromUrl() {
      this.error = ''
      this.parsed = null
      this.fetching = true
      try {
        const resp = await fetch(this.fetchUrl)
        if (!resp.ok) throw new Error(`HTTP ${resp.status}`)
        const block = this.extractRouteBlock(await resp.json())
        if (!block) this.error = this.$t('rule.import.errNoArraysFetched')
        else this.setParsed(block)
      } catch (e: any) {
        this.error = this.$t('rule.import.errFetch', { message: e.message })
      } finally {
        this.fetching = false
      }
    },
    async onFileUpload(files: File | File[] | null) {
      this.error = ''
      this.parsed = null
      const file = Array.isArray(files) ? files[0] : files
      if (!file) {
        this.error = this.$t('rule.import.errNoFile')
        return
      }
      try {
        const block = this.extractRouteBlock(JSON.parse(await file.text()))
        if (!block) {
          this.error = this.$t('rule.import.errNoArraysInFile')
          return
        }
        this.setParsed(block)
      } catch (e: any) {
        this.error = this.$t('rule.import.errJsonParse', { message: e.message })
        return
      }
    },
    save() {
      if (!this.parsed) return
      this.$emit('save', this.parsed, this.mode, this.applyFinal)
    },
    close() {
      this.$emit('close')
    },
  },
  watch: {
    visible(v: boolean) {
      if (v) this.reset()
    },
  },
}
</script>
