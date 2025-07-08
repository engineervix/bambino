<template>
  <v-container>
    <v-row>
      <v-col cols="12">
        <h1 class="text-h4 mb-4">Account</h1>
        
        <!-- User info -->
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

        <!-- Baby profiles -->
        <v-card class="mb-4">
          <v-card-title>Baby Profiles</v-card-title>
          <v-card-text>
            <v-list v-if="babies.length > 0">
              <v-list-item
                v-for="baby in babies"
                :key="baby.id"
                :active="baby.id === currentBaby?.id"
                @click="selectBaby(baby)"
                class="rounded"
              >
                <template v-slot:prepend>
                  <v-icon>mdi-baby-face</v-icon>
                </template>
                <v-list-item-title>{{ baby.name }}</v-list-item-title>
                <v-list-item-subtitle>
                  Born {{ formatDate(baby.birth_date) }} • {{ baby.age_display }}
                </v-list-item-subtitle>
                <template v-slot:append v-if="baby.id === currentBaby?.id">
                  <v-icon color="primary">mdi-check-circle</v-icon>
                </template>
              </v-list-item>
            </v-list>
            <div v-else class="text-center py-4">
              <p class="text-grey mb-2">No baby profiles yet</p>
              <p class="text-caption text-grey">
                Use the command line to create a baby profile:<br>
                <code>./baby-tracker create-user -u username -b "Baby Name"</code>
              </p>
            </div>
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
import { format } from 'date-fns'

const authStore = useAuthStore()
const { username, loading, babies, currentBaby } = storeToRefs(authStore)
const { logout, selectBaby } = authStore

async function handleLogout() {
  await logout()
}

function formatDate(dateString) {
  return format(new Date(dateString), 'MMM d, yyyy')
}
</script>