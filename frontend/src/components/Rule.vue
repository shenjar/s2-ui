<template>
  <ExpTextarea
    v-model="expTextarea.visible"
    :visible="expTextarea.visible"
    :label="expTextarea.title"
    :content="expTextarea.content"
    @update="saveExpTextarea"
    @close="closeExpTextarea"
  />
  <v-card style="background-color: inherit;">
    <v-row>
      <v-col cols="12" v-if="optionInbound">
        <v-combobox
          v-model="rule.inbound"
          :items="inTags"
          :label="$t('pages.inbounds')"
          multiple
          chips
          hide-details
        ></v-combobox>
      </v-col>
      <v-col cols="12" v-if="optionClient">
        <v-combobox
          v-model="rule.auth_user"
          :items="clients"
          :label="$t('pages.clients')"
          multiple
          chips
          hide-details
        ></v-combobox>
      </v-col>
      <v-col cols="12" sm="6" md="4" v-if="optionIPver">
        <v-select
          hide-details
          :label="$t('rule.ipVer')"
          :items="[4,6]"
          v-model.number="rule.ip_version">
        </v-select>
      </v-col>
      <v-col cols="12" sm="6" md="4" v-if="optionNetwork">
        <v-select
          hide-details
          multiple
          chips
          :label="$t('network')"
          :items="['tcp','udp','icmp']"
          v-model="rule.network">
        </v-select>
      </v-col>
      <v-col cols="12" sm="6" v-if="optionProtocol">
        <v-select
          v-model="rule.protocol"
          :items="protocols"
          :label="$t('protocol')"
          multiple
          chips
          hide-details
        ></v-select>
      </v-col>
    </v-row>
    <v-row v-if="optionDomain">
      <v-col cols="12" sm="6" md="4">
        <v-select
          hide-details
          :items="domainKeys"
          @update:model-value="updateDomainOption($event)"
          v-model="domainOption">
        </v-select>
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.domain != undefined">
        <v-textarea :label="$t('rule.domain')"
          hide-details
          v-model="domain"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.domain'), 'domain')"
        />
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.domain_suffix != undefined">
        <v-textarea :label="$t('rule.domainSufix')"
          hide-details
          v-model="domain_suffix"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.domainSufix'), 'domain_suffix')"
        />
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.domain_keyword != undefined">
        <v-textarea :label="$t('rule.domainKw')"
          hide-details
          v-model="domain_keyword"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.domainKw'), 'domain_keyword')"
        />
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.domain_regex != undefined">
        <v-textarea :label="$t('rule.domainRgx')"
          hide-details
          v-model="domain_regex"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.domainRgx'), 'domain_regex')"
        />
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.ip_cidr != undefined">
        <v-textarea :label="$t('rule.ip')"
          hide-details
          v-model="ip_cidr"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.ip'), 'ip_cidr')"
        />
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.ip_is_private != undefined">
        <v-switch v-model="rule.ip_is_private" color="primary" :label="$t('rule.privateIp')" hide-details></v-switch>
      </v-col>
    </v-row>
    <v-row v-if="optionPort">
      <v-col cols="12" sm="6" md="4">
        <v-select
          hide-details
          :items="portKeys"
          @update:model-value="updatePortOption($event)"
          v-model="portOption">
        </v-select>
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.port != undefined">
        <v-textarea :label="$t('rule.port')"
          hide-details
          v-model="port"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.port'), 'port')"
        />
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.port_range != undefined">
        <v-textarea :label="$t('rule.portRange')"
          hide-details
          v-model="port_range"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.portRange'), 'port_range')"
        />
      </v-col>
    </v-row>
    <v-row v-if="optionSrcIP">
      <v-col cols="12" sm="6" md="4">
        <v-select
          hide-details
          :items="srcIPKeys"
          @update:model-value="updateSrcIPOption($event)"
          v-model="srcIPOption">
        </v-select>
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.source_ip_cidr != undefined">
        <v-textarea :label="$t('rule.srcCidr')"
          hide-details
          v-model="source_ip_cidr"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.srcCidr'), 'source_ip_cidr')"
        />
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.source_ip_is_private != undefined">
        <v-switch v-model="rule.source_ip_is_private" color="primary" :label="$t('rule.srcPrivateIp')" hide-details></v-switch>
      </v-col>
    </v-row>
    <v-row v-if="optionSrcPort">
      <v-col cols="12" sm="6" md="4">
        <v-select
          hide-details
          :items="srcPortKeys"
          @update:model-value="updateSrcPortOption($event)"
          v-model="srcPortOption">
        </v-select>
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.source_port != undefined">
        <v-textarea :label="$t('rule.srcPort')"
          hide-details
          v-model="source_port"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.srcPort'), 'source_port')"
        />
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.source_port_range != undefined">
        <v-textarea :label="$t('rule.srcPortRange')"
          hide-details
          v-model="source_port_range"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.srcPortRange'), 'source_port_range')"
        />
      </v-col>
    </v-row>
    <v-row v-if="optionPreferredBy">
      <v-col cols="12" sm="6">
        <v-combobox
          v-model="rule.preferred_by"
          :items="outTags || inTags"
          :label="$t('rule.preferredBy')"
          multiple
          chips
          hide-details
        ></v-combobox>
      </v-col>
    </v-row>
    <v-row v-if="optionInterface">
      <v-col cols="12" sm="6" md="4">
        <v-select
          hide-details
          :items="interfaceKeys"
          @update:model-value="updateInterfaceOption($event)"
          v-model="interfaceOption">
        </v-select>
      </v-col>
      <v-col cols="12" sm="6" v-if="rule.interface_address != undefined || rule.network_interface_address != undefined || rule.default_interface_address != undefined">
        <v-textarea :label="$t('rule.interfaceAddr')"
          hide-details
          v-model="interface_addr"
          rows="5"
          no-resize
          density="compact"
          append-icon="mdi-arrow-expand"
          @click:append="openExpTextarea($t('rule.interfaceAddr'), 'interface_address')"
        />
      </v-col>
    </v-row>
    <v-row v-if="optionRuleSet">
      <v-col cols="12" sm="6">
        <v-combobox
          v-model="rule.rule_set"
          :items="rsTags"
          :label="$t('rule.ruleset')"
          multiple
          chips
          hide-details
        ></v-combobox>
      </v-col>
      <v-col cols="12" sm="6">
        <v-switch v-model="rule.rule_set_ip_cidr_match_source" color="primary" :label="$t('rule.rulesetMatchSrc')" hide-details></v-switch>
      </v-col>
    </v-row>
    <v-card-actions>
      <v-spacer></v-spacer>
      <v-menu v-model="menu" :close-on-content-click="false" location="start">
        <template v-slot:activator="{ props }">
          <v-btn v-bind="props" hide-details variant="tonal">{{ $t('rule.options') }}</v-btn>
        </template>
        <v-card>
          <v-list>
            <v-list-item>
              <v-switch v-model="optionInbound" color="primary" :label="$t('pages.inbounds')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionClient" color="primary" :label="$t('pages.clients')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionIPver" color="primary" :label="$t('rule.ipVer')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionNetwork" color="primary" :label="$t('network')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionProtocol" color="primary" :label="$t('protocol')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionDomain" color="primary" :label="$t('rule.domainRules')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionPort" color="primary" :label="$t('in.port')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionSrcIP" color="primary" :label="$t('rule.srcIpRules')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionSrcPort" color="primary" :label="$t('rule.srcPortRules')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionPreferredBy" color="primary" :label="$t('rule.preferredBy')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionInterface" color="primary" :label="$t('rule.interfaceAddr')" hide-details></v-switch>
            </v-list-item>
            <v-list-item>
              <v-switch v-model="optionRuleSet" color="primary" :label="$t('rule.ruleset')" hide-details></v-switch>
            </v-list-item>
          </v-list>
        </v-card>
      </v-menu>
    </v-card-actions>
  </v-card>
</template>

<script lang="ts">
import ExpTextarea from '@/components/ExpTextarea.vue'
export default {
  components: { ExpTextarea },
  props: ['rule', 'clients', 'inTags', 'outTags', 'rsTags', 'deleteable'],
  data() {
    return {
      menu: false,
      domainKeys: ['domain', 'domain_suffix', 'domain_keyword', 'domain_regex', 'ip_cidr', 'ip_is_private'],
      interfaceKeys: ['interface_address', 'network_interface_address', 'default_interface_address'],
      portKeys: ['port', 'port_range'],
      srcIPKeys: ['source_ip_cidr', 'source_ip_is_private'],
      srcPortKeys: ['source_port', 'source_port_range'],
      domainOption: 'domain',
      interfaceOption: 'interface_address',
      portOption: 'port',
      srcIPOption: 'source_ip_cidr',
      srcPortOption: 'source_port',
      protocols: [
        { title: 'HTTP', value: 'http' },
        { title: 'TLS', value: 'tls' },
        { title: 'QUIC', value: 'quic' },
        { title: 'STUN', value: 'stun' },
        { title: 'DNS', value: 'dns' },
        { title: 'BitTorrent', value: 'bittorrent' },
        { title: 'DTLS', value: 'dtls' },
        { title: 'SSH', value: 'ssh' },
        { title: 'RDP', value: 'rdp' },
        { title: 'NTP', value: 'ntp' },
      ],
      expTextarea: {
        visible: false,
        title: '',
        content: '',
        object: '',
      }
    }
  },
  methods: {
    updateDomainOption(option:string) {
      this.domainKeys.forEach(k => delete this.$props.rule[k])
      this.$props.rule[option] = option == 'ip_is_private' ? false : []
    },
    updatePortOption(option:string) {
      this.portKeys.forEach(k => delete this.$props.rule[k])
      this.$props.rule[option] = []
    },
    updateSrcIPOption(option:string) {
      this.srcIPKeys.forEach(k => delete this.$props.rule[k])
      this.$props.rule[option] = option == 'source_ip_is_private' ? false : []
    },
    updateSrcPortOption(option:string) {
      this.srcPortKeys.forEach(k => delete this.$props.rule[k])
      this.$props.rule[option] = []
    },
    updateInterfaceOption(option:string) {
      this.interfaceKeys.forEach(k => delete this.$props.rule[k])
      this.$props.rule[option] = []
    },
    openExpTextarea(title:string, object:string) {
      this.expTextarea.visible = !this.expTextarea.visible
      this.expTextarea.title = title
      this.expTextarea.content = this.$props.rule[object]?.join('\n') ?? ''
      this.expTextarea.object = object
    },
    saveExpTextarea(results:string[]) {
      this.$props.rule[this.expTextarea.object] = results
      this.closeExpTextarea()
    },
    closeExpTextarea() {
      this.expTextarea.visible = false
      this.expTextarea.title = ''
      this.expTextarea.object = ''
    },
  },
  computed: {
    optionInbound: {
      get() { return this.$props.rule.inbound != undefined },
      set(v:boolean) { this.$props.rule.inbound = v ? [] : undefined }
    },
    optionClient: {
      get() { return this.$props.rule.auth_user != undefined },
      set(v:boolean) { this.$props.rule.auth_user = v ? [] : undefined }
    },
    optionIPver: {
      get() { return this.$props.rule.ip_version != undefined },
      set(v:boolean) { this.$props.rule.ip_version = v ? 4 : undefined }
    },
    optionProtocol: {
      get() { return this.$props.rule.protocol != undefined },
      set(v:boolean) { this.$props.rule.protocol = v ? ['http'] : undefined }
    },
    optionDomain: {
      get() { return Object.keys(this.$props.rule).some(r => this.domainKeys.includes(r)) },
      set(v:boolean) { 
        if (v) {
          this.$props.rule.domain = []
        } else {
          this.domainKeys.forEach(k => delete this.$props.rule[k])
        }
        this.domainOption = 'domain'
      }
    },
    optionPort: {
      get() { return Object.keys(this.$props.rule).some(r => this.portKeys.includes(r)) },
      set(v:boolean) { 
        if (v) {
          this.$props.rule.port = []
        } else {
          this.portKeys.forEach(k => delete this.$props.rule[k])
        }
        this.portOption = 'port'
      }
    },
    optionSrcIP: {
      get() { return Object.keys(this.$props.rule).some(r => this.srcIPKeys.includes(r)) },
      set(v:boolean) { 
        if (v) {
          this.$props.rule.source_ip_cidr = []
        } else {
          this.srcIPKeys.forEach(k => delete this.$props.rule[k])
        }
        this.srcIPOption = 'source_ip_cidr'
      }
    },
    optionSrcPort: {
      get() { return Object.keys(this.$props.rule).some(r => this.srcPortKeys.includes(r)) },
      set(v:boolean) { 
        if (v) {
          this.$props.rule.source_port = []
        } else {
          this.srcPortKeys.forEach(k => delete this.$props.rule[k])
        }
        this.srcPortOption = 'source_port'
      }
    },
    optionPreferredBy: {
      get() { return this.$props.rule.preferred_by != undefined },
      set(v:boolean) { this.$props.rule.preferred_by = v ? [] : undefined }
    },
    optionInterface: {
      get() { return this.interfaceKeys.some(k => this.$props.rule[k] != undefined) },
      set(v:boolean) {
        if (v) {
          this.$props.rule.interface_address = []
        } else {
          this.interfaceKeys.forEach(k => delete this.$props.rule[k])
        }
        this.interfaceOption = 'interface_address'
      }
    },
    optionRuleSet: {
      get() { return this.$props.rule.rule_set != undefined },
      set(v:boolean) {
        if (v) {
          this.$props.rule.rule_set = []
          this.$props.rule.rule_set_ip_cidr_match_source = false
        } else {
          delete this.$props.rule.rule_set
          delete this.$props.rule.rule_set_ip_cidr_match_source
        }
      }
    },
    optionNetwork: {
      get() { return this.$props.rule.network != undefined },
      set(v:boolean) { this.$props.rule.network = v ? [] : undefined }
    },
    domain: {
      get() { return this.$props.rule.domain?.join('\n') ?? '' },
      set(v:string) { this.$props.rule.domain = v.length > 0 ? v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0) : [] }
    },
    domain_suffix: {
      get() { return this.$props.rule.domain_suffix?.join('\n') ?? '' },
      set(v:string) { this.$props.rule.domain_suffix = v.length > 0 ? v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0) : [] }
    },
    domain_keyword: {
      get() { return this.$props.rule.domain_keyword?.join('\n') ?? '' },
      set(v:string) { this.$props.rule.domain_keyword = v.length > 0 ? v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0) : [] }
    },
    domain_regex: {
      get() { return this.$props.rule.domain_regex?.join('\n') ?? '' },
      set(v:string) { this.$props.rule.domain_regex = v.length > 0 ? v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0) : [] }
    },
    ip_cidr: {
      get() { return this.$props.rule.ip_cidr?.join('\n') ?? '' },
      set(v:string) { this.$props.rule.ip_cidr = v.length > 0 ? v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0) : [] }
    },
    port: {
      get() { return this.$props.rule.port?.join('\n') ?? '' },
      set(v:string) {
        const lines = v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0)
        if (!v.endsWith('\n')) {
          this.$props.rule.port = lines.length > 0 ? lines.map((str:string) => parseInt(str, 10)).filter((n:number) => !isNaN(n)) : []
        }
      }
    },
    port_range: {
      get() { return this.$props.rule.port_range?.join('\n') ?? '' },
      set(v:string) { this.$props.rule.port_range = v.length > 0 ? v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0) : [] }
    },
    source_ip_cidr: {
      get() { return this.$props.rule.source_ip_cidr?.join('\n') ?? '' },
      set(v:string) { this.$props.rule.source_ip_cidr = v.length > 0 ? v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0) : [] }
    },
    source_port: {
      get() { return this.$props.rule.source_port?.join('\n') ?? '' },
      set(v:string) {
        const lines = v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0)
        if (!v.endsWith('\n')) {
          this.$props.rule.source_port = lines.length > 0 ? lines.map((str:string) => parseInt(str, 10)).filter((n:number) => !isNaN(n)) : []
        }
      }
    },
    source_port_range: {
      get() { return this.$props.rule.source_port_range?.join('\n') ?? '' },
      set(v:string) { this.$props.rule.source_port_range = v.length > 0 ? v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0) : [] }
    },
    interface_addr: {
      get() {
        const k = this.interfaceKeys.find(k => this.$props.rule[k] != undefined)
        return k ? this.$props.rule[k]?.join('\n') ?? '' : ''
      },
      set(v:string) {
        const k = this.interfaceKeys.find(k => this.$props.rule[k] != undefined)
        if (k) this.$props.rule[k] = v.length > 0 ? v.split('\n').map((s:string) => s.trim()).filter((s:string) => s.length > 0) : []
      }
    },
  },
  mounted() {
    const ruleKeys = Object.keys(this.$props.rule)
    if (this.optionDomain) {
      const enabledOption = this.domainKeys.filter(k => ruleKeys.includes(k))
      this.domainOption = enabledOption.length>0 ? enabledOption[0] : 'domain'
    }
    if (this.optionPort) {
      const enabledOption = this.portKeys.filter(k => ruleKeys.includes(k))
      this.portOption = enabledOption.length>0 ? enabledOption[0] : 'port'
    }
    if (this.optionSrcIP) {
      const enabledOption = this.srcIPKeys.filter(k => ruleKeys.includes(k))
      this.srcIPOption = enabledOption.length>0 ? enabledOption[0] : 'source_ip_cidr'
    }
    if (this.optionSrcPort) {
      const enabledOption = this.srcPortKeys.filter(k => ruleKeys.includes(k))
      this.srcPortOption = enabledOption.length>0 ? enabledOption[0] : 'source_port'
    }
    if (this.optionInterface) {
      const enabledOption = this.interfaceKeys.filter(k => ruleKeys.includes(k))
      this.interfaceOption = enabledOption.length>0 ? enabledOption[0] : 'interface_address'
    }
  }
}
</script>
