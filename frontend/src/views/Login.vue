<template>
  <v-app>
  <div class="login-page">
    <div class="login-blob login-blob-1"></div>
    <div class="login-blob login-blob-2"></div>
    <v-container class="fill-height login-container">
      <v-row justify="center" align="center" class="fill-height">
        <v-col cols="12" sm="8" md="5" lg="4" xl="3">
          <v-card class="login-card pa-2" elevation="10" rounded="xl">
            <v-card-text>
              <div class="text-center mb-6 mt-2">
                <v-img src="@/assets/logo.svg" width="72" height="72" class="mx-auto mb-3" />
                <div class="login-brand">2S-UI</div>
                <div class="text-medium-emphasis text-body-2 mt-1">{{ $t('login.title') }}</div>
              </div>
              <v-form @submit.prevent="login" ref="form">
                <v-text-field
                  v-model="username"
                  :label="$t('login.username')"
                  :rules="usernameRules"
                  prepend-inner-icon="mdi-account-outline"
                  autocomplete="username"
                  required
                  class="mb-2"
                ></v-text-field>
                <v-text-field
                  v-model="password"
                  :label="$t('login.password')"
                  :rules="passwordRules"
                  :type="showPassword ? 'text' : 'password'"
                  prepend-inner-icon="mdi-lock-outline"
                  :append-inner-icon="showPassword ? 'mdi-eye-off-outline' : 'mdi-eye-outline'"
                  @click:append-inner="showPassword = !showPassword"
                  autocomplete="current-password"
                  required
                ></v-text-field>
                <v-btn
                  :loading="loading"
                  type="submit"
                  color="primary"
                  block
                  size="large"
                  class="mt-3"
                  v-text="$t('actions.submit')"
                ></v-btn>
              </v-form>
              <v-divider class="my-5"></v-divider>
              <v-row align="center" no-gutters>
                <v-col>
                  <v-select
                    density="compact"
                    hide-details
                    variant="outlined"
                    prepend-inner-icon="mdi-translate"
                    :items="languages"
                    v-model="$i18n.locale"
                    @update:modelValue="changeLocale">
                  </v-select>
                </v-col>
                <v-col cols="auto" class="ps-3">
                  <v-btn-toggle density="compact" variant="outlined" divided rounded="lg" :model-value="currentTheme">
                    <v-btn
                      v-for="th in themes"
                      :key="th.value"
                      :value="th.value"
                      :icon="th.icon"
                      size="small"
                      v-tooltip:top="$t(`theme.${th.value}`)"
                      @click="changeTheme(th.value)"
                    ></v-btn>
                  </v-btn-toggle>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
  </v-app>
</template>

<script lang="ts" setup>
import { ref } from "vue"
import { useLocale, useTheme } from 'vuetify'
import { i18n, languages } from '@/locales'
import { useRouter } from 'vue-router'
import HttpUtil from '@/plugins/httputil'

const theme = useTheme()
const locale = useLocale()

const themes = [
  { value: 'light', icon: 'mdi-white-balance-sunny' },
  { value: 'dark', icon: 'mdi-moon-waning-crescent' },
  { value: 'system', icon: 'mdi-laptop' },
]

const currentTheme = ref(localStorage.getItem('theme') ?? 'system')

const username = ref('')
const usernameRules = [
  (value: string) => {
    if (value?.length > 0) return true
    return i18n.global.t('login.unRules')
  },
]

const password = ref('')
const passwordRules = [
  (value: string) => {
    if (value?.length > 0) return true
    return i18n.global.t('login.pwRules')
  },
]

const showPassword = ref(false)
const loading = ref(false)
const router = useRouter()

const login = async () => {
  if (username.value == '' || password.value == '') return
  loading.value=true
  const response = await HttpUtil.post('api/login',{user: username.value, pass: password.value})
  if(response.success){
    setTimeout(() => {
      loading.value=false
      router.push('/')
    }, 500)
  } else {
    loading.value=false
  }
}
const changeLocale = (l: any) => {
  locale.current.value = l ?? 'en'
  localStorage.setItem('locale', locale.current.value)
}
const changeTheme = (th: string) => {
  theme.change(th)
  localStorage.setItem('theme', th)
  currentTheme.value = th
}
</script>

<style scoped>
.login-page {
  position: relative;
  min-height: 100dvh;
  overflow: hidden;
  background:
    radial-gradient(1200px 600px at 85% -10%, rgba(99, 102, 241, 0.18), transparent 60%),
    radial-gradient(900px 500px at -10% 110%, rgba(14, 165, 233, 0.16), transparent 60%),
    rgb(var(--v-theme-background));
}
.login-container {
  position: relative;
  z-index: 1;
  min-height: 100dvh;
}
.login-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(90px);
  opacity: 0.45;
  pointer-events: none;
}
.login-blob-1 {
  width: 420px;
  height: 420px;
  top: -120px;
  inset-inline-end: -80px;
  background: rgba(99, 102, 241, 0.5);
}
.login-blob-2 {
  width: 380px;
  height: 380px;
  bottom: -140px;
  inset-inline-start: -100px;
  background: rgba(14, 165, 233, 0.4);
}
.login-card {
  backdrop-filter: blur(14px);
  background-color: rgba(var(--v-theme-surface), 0.86) !important;
}
.login-brand {
  font-size: 1.5rem;
  font-weight: 700;
  letter-spacing: 0.12em;
}
</style>
