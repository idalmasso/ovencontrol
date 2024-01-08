<template>
  <v-container>
    <v-row
      ><v-btn color="blue" @click="$router.push({ name: 'Home' })"
        >Indietro</v-btn
      ></v-row
    >
  </v-container>
  <ListProgramsPage @buttonClicked="buttonClickedHandler" />
</template>

<script setup>
import { defineEmits } from "vue";
import ListProgramsPage from "@/components/ListProgramsPage.vue";
const emits = defineEmits({ programSelected: String });
function buttonClickedHandler(name) {
  emits("programSelected", name);
}
</script>
