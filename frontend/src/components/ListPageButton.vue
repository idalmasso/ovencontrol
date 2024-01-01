<template>
  <div class="align-center text-center list-page-button">
    <v-btn
      height="300"
      width="300"
      :class="['d-flex bg-secondary justify-center']"
    >
    </v-btn>
    <div upper>{{ title }}</div>
  </div>
</template>
<style>
.list-page-button {
  max-width: fit-content;
}
</style>
<script setup>
defineProps(["title"]);
</script>
