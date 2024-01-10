<template>
  <v-container>
    <v-row justify="center" class="fill-height" align-content="center">
      <v-col v-for="listItem in listItems" :key="listItem.title">
        <ListPageButton
          :title="listItem.title"
          :icon="listItem.icon"
          :color="listItem.color"
          @click="$emit('buttonClicked', listItem.name)"
        ></ListPageButton>
      </v-col>
    </v-row>
  </v-container>
</template>
<script setup>
import { defineProps, defineEmits } from "vue";
import ListPageButton from "./ListPageButton.vue";
defineProps(["listItems"]);
defineEmits({ buttonClicked: String });
</script>
