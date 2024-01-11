<template>
  <div class="align-center text-center list-page-button">
    <v-btn
      :color="color"
      height="200"
      width="200"
      :class="['d-flex justify-center']"
    >
      <v-icon size="150">{{ icon }}</v-icon>
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
import { defineProps } from "vue";
defineProps({"title":String, "icon":String, color:{type: String, default:"blue"}});
</script>
