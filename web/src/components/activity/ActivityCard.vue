<template>
  <v-card
    :color="color"
    class="activity-card"
    rounded="xl"
    elevation="0"
    @click="$emit('click')"
  >
    <!-- Header with title -->
    <v-card-title class="pb-2">
      <span class="text-h5">{{ title }}</span>
    </v-card-title>

    <!-- Content area -->
    <v-card-text class="pa-4">
      <div class="activity-content">
        <v-icon :icon="icon" size="48" class="mb-3" />
        <p class="text-body-1">{{ description }}</p>
      </div>
    </v-card-text>

    <!-- Action button -->
    <div class="activity-action">
      <v-btn
        icon
        size="large"
        color="white"
        elevation="2"
        @click.stop="$emit('add')"
      >
        <v-icon>mdi-plus</v-icon>
      </v-btn>
    </div>
  </v-card>
</template>

<script setup>
defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  icon: {
    type: String,
    required: true
  },
  color: {
    type: String,
    default: 'primary'
  }
})

defineEmits(['click', 'add'])
</script>

<style scoped>
.activity-card {
  position: relative;
  min-height: 180px;
  cursor: pointer;
  transition: all 0.2s;
}

.activity-card:hover {
  transform: translateY(-2px);
}

.activity-content {
  text-align: center;
  padding: 20px 0;
}

.activity-action {
  position: absolute;
  bottom: 20px;
  right: 20px;
}

.activity-action .v-btn {
  background: rgba(76, 175, 80, 0.9) !important;
}
</style>