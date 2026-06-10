<template>
  <v-dialog v-model="dialog" max-width="620">
    <v-card>
      <v-card-title>{{ label }}</v-card-title>
      <v-divider />
      <v-card-text>
        <v-row>
          <v-col>{{ $t('rule.etaHint') }}</v-col>
        </v-row>
        <v-row>
          <v-col>
            <v-textarea
              v-model="localText"
              :label="label"
              variant="outlined"
              rows="16"
              :counter="$t('count')"
              persistent-counter
              :counter-value="(v: string) => v.split('\n').filter((l: string) => l.trim().length > 0).length"
              spellcheck="false"
            ></v-textarea>
          </v-col>
        </v-row>
      </v-card-text>
      <v-card-actions>
        <v-btn @click="resetChanges" color="error" variant="plain">{{ $t('reset') }}</v-btn>
        <v-spacer />
        <v-btn @click="closeModal" color="primary" variant="outlined">{{ $t('actions.close') }}</v-btn>
        <v-btn @click="saveChanges" color="primary" variant="tonal">{{ $t('actions.save') }}</v-btn>
      </v-card-actions>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">

export default {
  props: ['visible', 'label', 'content'],
  emits: ['update', 'close'],
  data() {
    return {
      dialog: false,
      localText: '',
    }
  },
  watch: {
    visible(v) {
      if (v) {
        this.localText = this.content
      }
    },
  },
  computed: {
  },
  methods: {
    saveChanges() {
      const unique = [
        ...new Set(
          this.localText
            .split('\n')
            .map((l: string) => l.trim())
            .filter((l: string) => l.length > 0)
        ),
      ]
      this.$emit('update', unique)
      this.dialog = false
    },
    resetChanges() {
      this.localText = this.content
    },
    closeModal() {
      this.$emit('close')
    },
  },
}
</script>
