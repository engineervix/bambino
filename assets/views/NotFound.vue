<template>
  <v-container class="py-16 text-center">
    <v-icon size="64" color="error" class="mb-4">mdi-alert-circle-outline</v-icon>
    <h1 class="text-h4 font-weight-bold mb-2">404 – Page Not Found</h1>
    <p class="text-subtitle-1 mb-6">Oops! The page you are looking for doesn’t exist.</p>
    <v-btn color="primary" to="/">
      <v-icon start>mdi-home</v-icon>
      Back to Home
    </v-btn>
  </v-container>
</template>

<script setup>
// Nothing needed here – purely presentational
</script> 