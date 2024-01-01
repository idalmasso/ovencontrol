<template>
  <div class="align-center text-center list-page-button">
    <v-btn
      color="blue"
      height="300"
      width="300"
      :class="['d-flex justify-center']"
    >
      <v-icon size="250">{{ button.icon }}</v-icon>
    </v-btn>
    <div upper>{{ button.title }}</div>
  </div>
</template>
<style>
.list-page-button {
  max-width: fit-content;
}
</style>
<script setup>
defineProps(["button"]);
</script>
