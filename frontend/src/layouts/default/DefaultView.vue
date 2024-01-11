<template>
  <v-main>
    <div class="text-center">
      <h1>{{ $route.meta.title }}<span v-if="$route.params.programName"> - {{ $route.params.programName }}</span></h1> 
    </div>
    <router-view />
  </v-main>
</template>

<script setup>
//
</script>
