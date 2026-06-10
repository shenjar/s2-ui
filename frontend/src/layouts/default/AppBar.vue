<template>
  <v-app-bar flat class="app-bar" height="64">
    <v-app-bar-nav-icon v-if="isMobile" @click="$emit('toggleDrawer')" />
    <span v-else style="width: 16px"></span>
    <v-app-bar-title class="app-bar-title">{{ $t(<string>route.name) }}</v-app-bar-title>
    <v-menu>
      <template v-slot:activator="{ props }">
        <v-btn icon v-bind="props">
          <v-icon>mdi-translate</v-icon>
        </v-btn>
      </template>
      <v-list density="compact" nav>
        <v-list-item
          v-for="lang in languages"
          :key="lang.value"
          rounded="lg"
          color="primary"
          @click="changeLocale(lang.value)"
          :active="isActiveLocale(lang.value)"
        >
          <v-list-item-title>{{ lang.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
    <v-menu>
      <template v-slot:activator="{ props }">
        <v-btn icon v-bind="props" class="me-2">
          <v-icon>mdi-theme-light-dark</v-icon>
        </v-btn>
      </template>
      <v-list density="compact" nav>
        <v-list-item
          v-for="th in themes"
          :key="th.value"
          rounded="lg"
          color="primary"
          @click="changeTheme(th.value)"
          :prepend-icon="th.icon"
          :active="isActiveTheme(th.value)"
        >
          <v-list-item-title>{{ $t(`theme.${th.value}`) }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
  </v-app-bar>
</template>

<script lang="ts" setup>
import { useLocale, useTheme } from 'vuetify'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { languages } from '@/locales'

defineProps(['isMobile'])

const route = useRoute()
const { locale: i18nLocale } = useI18n()
const vuetifyLocale = useLocale()
const theme = useTheme()

const changeLocale = (l: string) => {
  i18nLocale.value = l
  vuetifyLocale.current.value = l
  localStorage.setItem('locale', l)
  window.location.reload()
}
const isActiveLocale = (l: string) => i18nLocale.value === l
const themes = [
  { value: 'light', icon: 'mdi-white-balance-sunny' },
  { value: 'dark', icon: 'mdi-moon-waning-crescent' },
  { value: 'system', icon: 'mdi-laptop' },
]

const changeTheme = (th: string) => {
  theme.change(th)
  localStorage.setItem('theme', th)
}
const isActiveTheme = (th: string) => {
  const current = localStorage.getItem('theme') ?? 'system'
  return current == th
}
</script>

<style>
.app-bar {
  border-bottom: 1px solid rgba(var(--v-border-color), var(--v-border-opacity));
  background-color: rgba(var(--v-theme-background), 0.82) !important;
  backdrop-filter: blur(12px);
}
.app-bar-title {
  font-weight: 600;
  font-size: 1.05rem;
}
</style>
