<template>
  <v-container class="py-6">
    <v-row justify="center">
      <v-col cols="12" md="10" lg="8">
        <h2 class="text-h5 font-weight-medium mb-4">Account</h2>

        <!-- Side-by-side cards on desktop -->
        <v-row dense>
          <!-- User info -->
          <v-col cols="12" md="4">
            <v-card variant="outlined" rounded="lg" class="mb-6">
              <v-card-text>
                <v-list density="comfortable" lines="three">
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
          </v-col>

          <!-- Baby profiles -->
          <v-col cols="12" md="8">
            <v-card variant="outlined" rounded="lg" class="mb-6">
              <v-card-title>Baby Profiles</v-card-title>
              <v-card-text>
                <v-list v-if="babies.length > 0" density="comfortable">
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
                    <code>./bambino create-user -u username -b "Baby Name"</code>
                  </p>
                </div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>

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
