<template>
  <v-container>
   <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn
      ></v-row
    >
    <OvenConfigPage />
  </v-container>
  </template>
  
  <script setup>
  import OvenConfigPage from "@/components/OvenConfigPage.vue";
  </script>
  