<template>
  <v-app-bar :elevation="14" color="blue" dark>
    <v-app-bar-nav-icon @click="$router.push({ name: 'Home' })">
      <img src="@/assets/logo.png" height="40px" contain
    /></v-app-bar-nav-icon>
    <v-app-bar-title>
      <router-link style="text-decoration: none; color: inherit" to="/">
        Oven Controller</router-link
      >
    </v-app-bar-title>
  </v-app-bar>
</template>

<script setup>
//
</script>
