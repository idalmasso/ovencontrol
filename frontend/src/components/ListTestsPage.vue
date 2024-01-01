<template>
  <v-container>
    <v-row justify="center" class="fill-height" align-content="center">
      <v-col>
        <ListPageButton
          title="Test temperatura"
          @click="$router.push({ name: 'TemperatureCheck' })"
        ></ListPageButton>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import ListPageButton from "./ListPageButton.vue";
</script>
