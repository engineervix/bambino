<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-4">Account</h1>
        
        <v-card class="mb-4">
          <v-card-text>
            <v-list>
              <v-list-item>
                <template v-slot:prepend>
                  <v-icon>mdi-account</v-icon>
                </template>
                <v-list-item-title>{{ username }}</v-list-item-title>
                <v-list-item-subtitle>Logged in user</v-list-item-subtitle>
              </v-list-item>
            </v-list>
          </v-card-text>
        </v-card>

        <v-btn
          color="error"
          size="large"
          block
          @click="handleLogout"
          :loading="loading"
        >
          <v-icon start>mdi-logout</v-icon>
          Sign Out
        </v-btn>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { useAuthStore } from '@/stores/auth'
import { storeToRefs } from 'pinia'

const authStore = useAuthStore()
const { username, loading } = storeToRefs(authStore)
const { logout } = authStore

async function handleLogout() {
  await logout()
}
</script>