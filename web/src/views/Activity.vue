<template>
  <v-container class="pa-0">
    <!-- Header with baby info -->
    <v-sheet class="pa-4 mb-4" color="surface">
      <div class="d-flex align-center">
        <v-avatar size="56" class="mr-3">
          <v-icon size="large">mdi-baby-face</v-icon>
        </v-avatar>
        <div>
          <h1 class="text-h5">Baby</h1>
          <p class="text-body-2 text-grey">{{ currentDate }}</p>
        </div>
      </div>
    </v-sheet>

    <!-- Activity cards -->
    <v-container>
      <v-row>
        <v-col cols="12">
          <p class="text-center text-h6 mb-4">Quick Add Activities</p>
        </v-col>
      </v-row>
    </v-container>
  </v-container>
</template>

<script setup>
import { ref, computed } from 'vue'
import { format } from 'date-fns'

// Current date display
const currentDate = computed(() => {
  return format(new Date(), 'EEEE, MMM d')
})
</script>