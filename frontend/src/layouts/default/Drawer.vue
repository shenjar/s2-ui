<template>
  <v-navigation-drawer
    v-model="showDrawer"
    :temporary="isMobile"
    :expand-on-hover="!isMobile"
    :rail="!isMobile"
    :permanent="!isMobile"
    class="app-drawer"
    @click="isMobile ? $emit('toggleDrawer') : null"
  >
    <v-list-item
      height="64"
      class="app-drawer-brand"
      prepend-avatar="@/assets/logo.svg"
    >
      <v-list-item-title class="app-drawer-title">2S-UI</v-list-item-title>
      <template v-slot:append v-if="isMobile">
        <v-icon icon="mdi-close" />
      </template>
    </v-list-item>

    <v-divider></v-divider>

    <v-list density="compact" nav color="primary" class="app-drawer-nav">
      <template v-for="(group, gi) in menuGroups" :key="gi">
        <v-divider v-if="gi > 0" class="my-2 mx-3" />
        <v-list-item link
          v-for="item in group"
          :key="item.title"
          :to="item.path"
          rounded="lg"
          :active="router.currentRoute.value.path == item.path">
          <template v-slot:prepend>
            <v-icon :icon="item.icon" size="20"></v-icon>
          </template>
          <v-list-item-title v-text="$t(item.title)"></v-list-item-title>
        </v-list-item>
      </template>
    </v-list>
    <template v-slot:append>
      <v-divider></v-divider>
      <v-list density="compact" nav>
        <v-list-item rounded="lg" base-color="error" @click="Logout">
          <template v-slot:prepend>
            <v-icon icon="mdi-logout" size="20"></v-icon>
          </template>
          <v-list-item-title v-text="$t('menu.logout')"></v-list-item-title>
        </v-list-item>
      </v-list>
    </template>
  </v-navigation-drawer>
</template>

<script lang="ts" setup>
import { computed } from 'vue'
import router from '@/router'
import { logout } from '@/plugins/httputil'

const props = defineProps(['isMobile','displayDrawer'])

const showDrawer = computed((): boolean => {
  return props.displayDrawer
})

const menuGroups = [
  [
    { title: 'pages.home', icon: 'mdi-view-dashboard-outline', path: '/' },
  ],
  [
    { title: 'pages.inbounds', icon: 'mdi-cloud-download-outline', path: '/inbounds' },
    { title: 'pages.clients', icon: 'mdi-account-multiple-outline', path: '/clients' },
    { title: 'pages.outbounds', icon: 'mdi-cloud-upload-outline', path: '/outbounds' },
    { title: 'pages.endpoints', icon: 'mdi-cloud-tags', path: '/endpoints' },
    { title: 'pages.services', icon: 'mdi-server-outline', path: '/services' },
  ],
  [
    { title: 'pages.rules', icon: 'mdi-routes', path: '/rules' },
    { title: 'pages.dns', icon: 'mdi-dns-outline', path: '/dns' },
    { title: 'pages.tls', icon: 'mdi-certificate-outline', path: '/tls' },
    { title: 'pages.basics', icon: 'mdi-application-cog-outline', path: '/basics' },
  ],
  [
    { title: 'pages.admins', icon: 'mdi-shield-account-outline', path: '/admins' },
    { title: 'pages.settings', icon: 'mdi-cog-outline', path: '/settings' },
  ],
]

const Logout = async () => {
  logout()
}
</script>

<style>
.app-drawer .v-navigation-drawer__content {
  scrollbar-width: none;
}
.app-drawer-brand .v-avatar {
  border-radius: 0;
}
.app-drawer-title {
  font-weight: 700;
  letter-spacing: 0.08em;
}
.app-drawer-nav .v-list-item--active {
  font-weight: 600;
}
</style>
