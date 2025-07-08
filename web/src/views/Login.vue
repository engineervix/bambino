<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="pa-4">
          <v-card-title class="text-h4 text-center mb-4">
            <v-icon icon="mdi-baby-face" size="large" class="mr-2" />
            Baby Tracker
          </v-card-title>
          
          <v-card-subtitle class="text-center mb-6">
            Sign in to continue
          </v-card-subtitle>
          
          <v-card-text>
            <v-form @submit.prevent="handleLogin" ref="form">
              <v-text-field
                v-model="credentials.username"
                label="Username"
                variant="outlined"
                prepend-inner-icon="mdi-account"
                :rules="[v => !!v || 'Username is required']"
                :disabled="loading"
                class="mb-4"
              />
              
              <v-text-field
                v-model="credentials.password"
                label="Password"
                type="password"
                variant="outlined"
                prepend-inner-icon="mdi-lock"
                :rules="[v => !!v || 'Password is required']"
                :disabled="loading"
                @keyup.enter="handleLogin"
              />
              
              <v-alert
                v-if="error"
                type="error"
                variant="tonal"
                class="mb-4"
                closable
                @click:close="clearError"
              >
                {{ error }}
              </v-alert>
              
              <v-btn
                type="submit"
                color="primary"
                size="large"
                block
                :loading="loading"
                :disabled="loading"
              >
                Sign In
              </v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'
import { useRouter } from 'vue-router'

const router = useRouter()
const authStore = useAuthStore()
const { loading, error } = storeToRefs(authStore)
const { login, clearError, checkAuth } = authStore

const form = ref(null)
const credentials = ref({
  username: '',
  password: ''
})

// Check if already authenticated on mount
onMounted(async () => {
  const isAuth = await checkAuth()
  if (isAuth) {
    // Already logged in, redirect to home
    router.push('/')
  }
})

async function handleLogin() {
  const { valid } = await form.value.validate()
  
  if (valid) {
    await login(credentials.value)
  }
}
</script>

<style scoped>
.v-container {
  background: radial-gradient(circle at center, rgba(76, 175, 80, 0.1) 0%, transparent 70%);
}
</style>