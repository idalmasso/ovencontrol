<template>
  <v-container>
    <v-row justify="center" class="fill-height" align-content="center">
      <v-col v-for="listItem in listItems" :key="listItem.title">
        <ListPageButton
          :button="listItem.button"
          @click="listItem.action"
        ></ListPageButton>
      </v-col>
    </v-row>
  </v-container>
</template>
<script setup>
import ListPageButton from "./ListPageButton.vue";
defineProps(["listItems"]);
</script>
