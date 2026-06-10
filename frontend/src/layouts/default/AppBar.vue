<template>
  <v-app-bar :elevation="5">
    <v-icon v-if="isMobile" icon="mdi-menu" @click="$emit('toggleDrawer')" />
    <span v-else style="width: 24px"></span>
    <v-app-bar-title :text="$t(<string>route.name)" class="align-center text-center " />
    <v-menu>
      <template v-slot:activator="{ props }">
        <v-btn icon v-bind="props">
          <v-icon>mdi-translate</v-icon>
        </v-btn>
      </template>
      <v-list>
        <v-list-item
          v-for="lang in languages"
          :key="lang.value"
          @click="changeLocale(lang.value)"
          :active="isActiveLocale(lang.value)"
        >
          <v-list-item-title>{{ lang.title }}</v-list-item-title>
        </v-list-item>
      </v-list>
    </v-menu>
    <v-menu>
      <template v-slot:activator="{ props }">
        <v-btn icon v-bind="props">
          <v-icon>mdi-theme-light-dark</v-icon>
        </v-btn>
      </template>
      <v-list>
        <v-list-item
          v-for="th in themes"
          :key="th.value"
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
