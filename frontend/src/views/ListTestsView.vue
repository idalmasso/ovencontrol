<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn
      ></v-row
    >
  </v-container>
  <ListTestsPage />
</template>

<script setup>
import ListTestsPage from "@/components/ListTestsPage.vue";
</script>
