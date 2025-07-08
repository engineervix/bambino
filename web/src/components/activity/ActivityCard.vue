<template>
  <v-card
    :color="color"
    class="activity-card"
    rounded="lg"
    elevation="2"
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
        color="primary"
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
  transition: transform 0.25s ease, box-shadow 0.25s ease;
  border-radius: 12px;
  overflow: hidden;
  /* Subtle outline */
  border: 1px solid rgba(255, 255, 255, 0.14);
}

/* Light diagonal gradient wash */
.activity-card::before {
  content: "";
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.12) 0%, rgba(255, 255, 255, 0) 70%);
  pointer-events: none;
}

/* Subtle radial highlight that fades in on hover */
.activity-card::after {
  content: "";
  position: absolute;
  inset: 0;
  background: radial-gradient(circle at top right, rgba(255, 255, 255, 0.12), transparent 60%);
  pointer-events: none;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.activity-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 6px 16px rgba(0, 0, 0, 0.2);
}

.activity-card:hover::after {
  opacity: 1;
}

.activity-content {
  text-align: center;
  padding: 30px 12px 60px;
}

.activity-content .v-icon {
  font-size: 56px !important;
  margin-bottom: 16px;
}

.activity-action {
  position: absolute;
  bottom: 24px;
  right: 24px;
}

.activity-action .v-btn {
  transition: transform 0.25s ease;
}

.activity-action .v-btn:hover {
  transform: rotate(90deg) scale(1.15);
}
</style>