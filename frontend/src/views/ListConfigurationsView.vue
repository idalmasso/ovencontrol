<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn
      ></v-row
    >
  </v-container>
  <ListConfigurationsPage />
</template>

<script setup>
import ListConfigurationsPage from "@/components/ListConfigurationsPage.vue";
</script>
