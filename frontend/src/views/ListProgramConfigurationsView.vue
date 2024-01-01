<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn
      ></v-row
    >
  </v-container>
  <ListProgramConfigurationsPage />
</template>

<script setup>
import ListProgramConfigurationsPage from "@/components/ListProgramConfigurationsPage.vue";
</script>
